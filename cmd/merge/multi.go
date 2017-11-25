package merge

import (
	"fmt"

	"github.com/evovetech/got/git"
	"github.com/evovetech/got/git/merge"
	"github.com/evovetech/got/util"
)

type multi Merger

type branchRef struct {
	Name string
	Ref  git.Ref
}

func getBranchRef(ours git.Ref, theirs git.Ref) (*branchRef, error) {
	name := fmt.Sprintf("%s_merge_%s", ours.ShortName(), theirs.ShortSha())
	if err := git.Command("branch", name, ours.ShortSha()).Run(); err != nil {
		return nil, err
	}
	if ref, err := git.ParseRef(name); err != nil {
		return nil, err
	} else {
		return &branchRef{name, ref}, nil
	}
}

func (b *branchRef) String() string {
	return util.String(b)
}

func (b *branchRef) Update() error {
	return b.Ref.Update()
}

type multiStep struct {
	multi *multi

	ours     branchRef
	theirs   branchRef
	strategy merge.Strategy
}

func (m *multiStep) Update(st merge.Strategy) error {
	if err := m.ours.Update(); err != nil {
		return err
	}
	if err := m.theirs.Update(); err != nil {
		return err
	}
	m.strategy = st
	return nil
}

func (m *multi) Run() (RunStep, error) {
	var err error
	var ours, theirs *branchRef
	if ours, err = getBranchRef(m.HeadRef, m.MergeRef); err != nil {
		return nil, err
	}
	if theirs, err = getBranchRef(m.MergeRef, m.HeadRef); err != nil {
		return nil, err
	}
	fmt.Printf("ours: %s\n", ours)
	fmt.Printf("theirs: %s\n", theirs)
	reset := func(err error) (RunStep, error) {
		git.Merge().Abort()
		m.HeadRef.Checkout()
		m.HeadRef.Reset().Hard()
		ours.Ref.Delete()
		theirs.Ref.Delete()
		return nil, err
	}

	// first merge
	s1 := &multiStep{m, *ours, *theirs, m.Strategy}
	if err := s1.RunE(); err != nil {
		return reset(err)
	}

	// second merge
	if err := s1.Update(merge.THEIRS); err != nil {
		return reset(err)
	}
	if err := s1.RunE(); err != nil {
		return reset(err)
	}

	// final merge
	if err := s1.Update(merge.THEIRS); err != nil {
		return reset(err)
	}
	if err := s1.RunLast(); err != nil {
		return reset(err)
	}

	// update refs
	if err := m.HeadRef.Checkout(); err != nil {
		return reset(err)
	}
	mrge := git.Merge()
	mrge.MergeRef = s1.ours.Name
	mrge.FFOnly()
	if err = mrge.Run(); err != nil {
		return reset(err)
	}

	ours.Ref.Delete()
	theirs.Ref.Delete()

	return nil, nil //, m.RunE()
}

func (m *multiStep) RunE() error {
	st := m.strategy
	step1 := NewStep(m.ours, m.theirs.Ref, st)
	step2 := NewStep(m.theirs, m.ours.Ref, st)
	steps := NewStepper(step1, step2)
	return steps.RunE()
}

func (m *multiStep) RunLast() error {
	lastStep := NewStep(m.ours, m.theirs.Ref, merge.THEIRS)
	return lastStep.RunE()
}
