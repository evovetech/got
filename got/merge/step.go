/*
 * Copyright 2018 evove.tech
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
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
	"github.com/evovetech/got/git"
	"github.com/evovetech/got/got/merge/mv"
	"github.com/evovetech/got/got/resolve"
	"github.com/evovetech/got/log"
	"github.com/evovetech/got/util"
)

type Step struct {
	Branch      *branchRef
	Target      *branchRef
	FindRenames int
}

var (
	movesEnabled =false
)

func NewStep(branch *branchRef, target *branchRef, findRenames int) *Step {
	return &Step{
		Branch:      branch,
		Target:      target,
		FindRenames: findRenames,
	}
}

func (s *Step) mergeCommitTree(tree string) (commit string, err error) {
	ours := s.Branch
	p1 := s.Branch.Ref.ShortSha()
	p2 := s.Target.Ref.ShortSha()
	commitCmd := git.Command("commit-tree", tree,
		"-p", p1,
		"-p", p2,
		"-m", s.getMsg())
	if commit, err = commitCmd.Output(); err == nil {
		err = git.Command("update-ref", ours.Ref.Info.RefName, commit).Run()
	}
	return
}

func (s *Step) Run() error {
	// COPY original branches
	ours, theirs := *s.Branch, *s.Target
	cp := *s
	cp.Branch = &ours
	cp.Target = &theirs

	// try merge
	if err := s.merge(); err != nil {
		return err
	}

	if !movesEnabled {
		return s.commit()
	}

	if m, ok := mv.GetFileMoves(); ok {
		return s.commit()
	} else {
		if err := util.RunAll(
			m.Run,
			s.merge,
			s.commit,
		); err != nil {
			return err
		}

		tree := git.RevParse("HEAD").Short() + "^{tree}"
		commit, err := cp.mergeCommitTree(tree)
		log.Printf("commit: %s", commit)
		return err
	}
}

func (s *Step) merge() error {
	if err := util.RunAll(s.checkout, s.updateBranchRef); err != nil {
		return err
	}

	// git merge --no-commit -X "$cmd" "$merge_commit"
	st := git.OURS
	m := s.Target.Ref.Merge()
	m.Strategy = st
	m.NoCommit()
	if s.FindRenames != 0 {
		m.FindRenames(s.FindRenames)
	}
	if err := m.Run(); err != nil {
		if err := resolve.Run(st); err != nil {
			git.AbortMerge()
			return fmt.Errorf("could not merge %s: %s", s.Target.Name, err.Error())
		}
	}
	return nil
}

func (s *Step) commit() error {
	return git.Command("commit", "-m", s.getMsg()).Run()
}

func (s *Step) checkout() error {
	return s.Branch.Ref.Checkout()
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

func getMsg(resolveName string) string {
	format := "resolving conflicts with %s changes"
	return fmt.Sprintf(format, resolveName)
}

func (s *Step) getMsg() string {
	head := s.Branch.OursName
	target := s.Target.OursName
	format := "merge %s into %s -- CONFLICTS -- resolving with %s changes"
	return fmt.Sprintf(format, target, head, head)
}
