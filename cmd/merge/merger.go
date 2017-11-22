/*
 * Copyright 2017 evove.tech
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package merge

import (
	"github.com/evovetech/got/git"
	"github.com/evovetech/got/git/merge"
	"github.com/evovetech/got/util"
)

type Merger struct {
	HeadRef  git.Ref
	MergeRef git.Ref
	Strategy merge.Strategy
}

func (m *Merger) String() string {
	return util.String(m)
}

func (m *Merger) RunE() error {
	if err := m.trySimpleMerge(); err == nil {
		// done
		return nil
	}

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

func (m *Merger) trySimpleMerge() error {
	simple := NewStep(m.HeadRef, m.MergeRef, merge.NONE)
	return simple.RunE()
}
