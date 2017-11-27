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
	"fmt"

	"github.com/evovetech/got/cmd/resolve"
	"github.com/evovetech/got/git"
	"github.com/evovetech/got/git/merge"
	"github.com/evovetech/got/util"
)

type Step struct {
	Branch   branchRef
	Target   branchRef
	Strategy merge.Strategy
}

func NewStep(branch branchRef, target branchRef, strategy merge.Strategy) *Step {
	return &Step{
		Branch:   branch,
		Target:   target,
		Strategy: strategy,
	}
}

func (s *Step) RunE() error {
	if err := util.RunAll(s.checkout, s.updateBranchRef); err != nil {
		return err
	}

	// git merge --no-commit -X "$cmd" "$merge_commit"
	m := s.Target.Ref.Merge()
	m.NoCommit()
	m.Strategy = s.Strategy
	if err := m.Run(); err != nil {
		if err := resolve.Run(s.Strategy); err != nil {
			git.Merge().Abort()
			return fmt.Errorf("could not merge %s: %s", s.Target.Name, err.Error())
		}
	}

	return s.commit()
}

func (s *Step) commit() error {
	return git.Command("commit", "-m", s.getMsg()).Run()
}

func (s *Step) checkout() error {
	if err := CheckStatus(); err != nil {
		return err
	}

	checkout := git.Checkout()
	checkout.AddArg(s.Branch.Name)
	return checkout.Run()
}

func (s *Step) updateBranchRef() error {
	// update ref and msg
	var err error
	var branchRef git.Ref
	if branchRef, err = git.ParseRef(s.Branch.Name); err != nil {
		return err
	}
	s.Branch.Ref = branchRef
	return nil
}

func (s *Step) deleteBranch() error {
	b := s.Branch.Name
	if b != "" {
		if err := git.Command("branch", "-D", b).Run(); err != nil {
			return err
		}
	}
	return nil
}

func (s *Step) getMsg() string {
	head := s.Branch.OrigName
	target := s.Target.OrigName
	format := "merge %s into %s -- CONFLICTS -- resolving with %s changes"
	resolve := target
	if s.Strategy == merge.OURS {
		resolve = head
	}
	return fmt.Sprintf(format, target, head, resolve)
}
