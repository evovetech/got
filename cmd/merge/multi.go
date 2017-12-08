package merge

import (
	"fmt"

	"github.com/evovetech/got/git"
	"github.com/evovetech/got/git/merge"
	"github.com/evovetech/got/util"
)

type multi Merger

type branchRef struct {
	OursName string
	Name     string
	Ref      git.Ref
}

func getBranchRef(ours git.Ref, theirs git.Ref) (*branchRef, error) {
	name := fmt.Sprintf("%s_merge_%s", ours.ShortName(), theirs.ShortSha())
	branchCmd := git.Command("branch", name, ours.ShortSha())
	if err := branchCmd.Run(); err != nil {
		if e := git.Command("branch", "-D", name).Run(); e != nil {
			return nil, err
		}
		if err := branchCmd.Run(); err != nil {
			return nil, err
		}
	}
	if ref, err := git.ParseRef(name); err != nil {
		return nil, err
	} else {
		return &branchRef{ours.ShortName(), name, ref}, nil
	}
}

func (b *branchRef) String() string {
	return util.String(b)
}

func (b *branchRef) Update() error {
	return b.Ref.Update()
}

func (b *branchRef) SetOursName(oursName string) {
	b.OursName = oursName
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
	reset := func(err error) (RunStep, error) {
		git.AbortMerge()
		m.HeadRef.Checkout()
		m.HeadRef.Reset().Hard()
		ours.Ref.Delete()
		theirs.Ref.Delete()
		return nil, err
	}

	// first merge
	s1 := &multiStep{m, *ours, *theirs, merge.OURS}
	if err := s1.RunE(25); err != nil {
		return reset(err)
	}

	// second merge
	if m.Strategy == merge.OURS {
		if err := s1.Update(merge.THEIRS); err != nil {
			return reset(err)
		}
		if err := s1.MergeResetTheirs(); err != nil {
			return reset(err)
		}
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

func (m *multiStep) MergeResetTheirs() error {
	st := merge.THEIRS
	err := util.RunAll(
		NewStep(m.ours, m.theirs, st).MergeResetTheirs,
		NewStep(m.theirs, m.ours, st).MergeResetTheirs,
	)
	if err != nil {
		return err
	}
	oursName := m.theirs.OursName
	m.theirs.SetOursName(m.ours.OursName)
	m.ours.SetOursName(oursName)
	return nil
}

func (m *multiStep) RunE(findRenames int) error {
	st := m.strategy
	step1 := NewStep(m.ours, m.theirs, st)
	step1.FindRenames = findRenames
	step2 := NewStep(m.theirs, m.ours, st)
	step2.FindRenames = findRenames
	steps := NewStepper(step1, step2)
	if err := steps.Run(); err != nil {
		return err
	}
	if st == merge.THEIRS {
		oursName := m.theirs.OursName
		m.theirs.SetOursName(m.ours.OursName)
		m.ours.SetOursName(oursName)
	}
	return nil
}

func (m *multiStep) RunLast() error {
	lastStep := NewStep(m.ours, m.theirs, merge.THEIRS)
	return lastStep.MergeResetTheirs()
}
