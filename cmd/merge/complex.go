package merge

import (
	"github.com/evovetech/got/git"
	"github.com/evovetech/got/git/merge"
)

type multi Merger

func (m *multi) Run() (RunStep, error) {
	return nil, m.RunE()
}

func (m *multi) RunE() error {
	var cleanupSteps []*Step
	cleanup := func() {
		for _, s := range cleanupSteps {
			s.deleteBranch()
		}
	}

	// multi-step merge
	st := m.Strategy
	step1 := NewStep(m.HeadRef, m.MergeRef, st)
	step2 := NewStep(m.MergeRef, m.HeadRef, st)
	steps := NewStepper(step1, step2)
	if err := steps.RunE(); err != nil {
		return err
	}

	cleanupSteps = append(cleanupSteps, step1, step2)
	headMergeRef, err := git.ParseRef(step1.branch)
	if err != nil {
		step1.reset()
		cleanup()
		return err
	}
	mergeRef, err := git.ParseRef(step2.branch)
	if err != nil {
		step1.reset()
		cleanup()
		return err
	}
	step3 := NewStep(headMergeRef, mergeRef, merge.THEIRS)
	step4 := NewStep(mergeRef, headMergeRef, merge.THEIRS)
	steps = NewStepper(step3, step4)
	if err := steps.RunE(); err != nil {
		step1.reset()
		cleanup()
		return err
	}

	cleanupSteps = append(cleanupSteps, step3, step4)
	oursRef, err := git.ParseRef(step3.branch)
	if err != nil {
		step1.reset()
		cleanup()
		return err
	}
	theirsRef, err := git.ParseRef(step4.branch)
	if err != nil {
		step1.reset()
		cleanup()
		return err
	}

	lastStep := NewStep(oursRef, theirsRef, merge.THEIRS)
	if err := lastStep.RunE(); err != nil {
		step1.reset()
		cleanup()
		return err
	}

	cleanupSteps = append(cleanupSteps, lastStep)
	if err := m.HeadRef.Checkout(); err != nil {
		step1.reset()
		cleanup()
		return err
	}

	mrge := git.Merge()
	mrge.MergeRef = lastStep.branch
	mrge.FFOnly()

	if err = mrge.Run(); err != nil {
		step1.reset()
	}
	cleanup()
	return err
}
