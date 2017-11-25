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
	"regexp"

	"github.com/evovetech/got/git"
	"github.com/evovetech/got/git/merge"
	"github.com/evovetech/got/util"
)

type Step struct {
	Branch      branchRef
	MergeTarget git.Ref
	Strategy    merge.Strategy
}

func NewStep(branch branchRef, mergeTarget git.Ref, strategy merge.Strategy) *Step {
	return &Step{
		Branch:      branch,
		MergeTarget: mergeTarget,
		Strategy:    strategy,
	}
}

func (s *Step) RunE() error {
	if err := util.RunAll(s.checkout, s.updateBranchRef); err != nil {
		return err
	}

	// git merge --no-commit -X "$cmd" "$merge_commit"
	m := s.MergeTarget.Merge()
	m.NoCommit()
	m.Strategy = s.Strategy
	if err := m.Run(); err != nil {
		if err := s.resolveUnmerged(); err != nil {
			git.Merge().Abort()
			return fmt.Errorf("could not merge %s: %s", s.MergeTarget.ShortName(), err.Error())
		}
	}

	return s.commit()
}

var reDD = regexp.MustCompile("^(DD)")
var reDeletedOurs = regexp.MustCompile("^(D|UA)")
var reDeletedTheirs = regexp.MustCompile("^(.D|AU)")

func (s *Step) resolveUnmerged() error {
	//git diff --name-only --diff-filter=UXB
	var errors []error
	diff := git.Command("diff", "--name-only", "--diff-filter=UXB")
	for _, file := range diff.OutputLines() {
		if err := s.resolveFile(file); err != nil {
			errors = append(errors, err)
		}
	}

	if len(errors) > 0 {
		var errString string
		for _, err := range errors {
			errString += fmt.Sprintln(err.Error())
		}
		return fmt.Errorf("%s", errString)
	}
	return nil
}

func (s *Step) resolveFile(file string) error {
	var err error
	st := s.Strategy
	status, err := git.Status(file).Output()
	if err != nil {
		return err
	}
	switch {
	case reDD.MatchString(status):
		err = git.Add(file).Run()
	case st == merge.OURS:
		err = s.resolveOurs(file, status)
	case st == merge.THEIRS:
		err = s.resolveTheirs(file, status)
	default:
		err = fmt.Errorf("unknown strategy: ")
	}
	return err
}

func (s *Step) resolveOurs(file string, status string) error {
	switch {
	case reDeletedOurs.MatchString(status):
		return git.ResolveRm(file).Run()
	default:
		return git.ResolveCheckout(file, merge.OURS).Run()
	}
}

func (s *Step) resolveTheirs(file string, status string) error {
	switch {
	case reDeletedTheirs.MatchString(status):
		return git.ResolveRm(file).Run()
	default:
		return git.ResolveCheckout(file, merge.THEIRS).Run()
	}
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
	head := s.Branch.Name
	target := s.MergeTarget.ShortName()
	format := "merge %s into %s -- CONFLICTS -- resolving with %s changes"
	resolve := target
	if s.Strategy == merge.OURS {
		resolve = head
	}
	return fmt.Sprintf(format, target, head, resolve)
}
