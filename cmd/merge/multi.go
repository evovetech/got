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
	"github.com/evovetech/got/git/merge"
	"github.com/evovetech/got/log"
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

func (b *branchRef) Update() error {
	return b.Ref.Update()
}

func (b branchRef) String() string {
	return util.String(b)
}

type multiStep struct {
	multi *multi

	ours   *branchRef
	theirs *branchRef
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
		log.Err.Println("aborting")
		git.AbortMerge()
		m.HeadRef.Checkout()
		m.HeadRef.Reset().Hard()
		ours.Ref.Delete()
		theirs.Ref.Delete()
		return nil, err
	}
	defer func() {
		if r := recover(); r != nil {
			log.Err.Println("recovering")
			reset(nil)
			panic(r)
		}
	}()
	if err = m.run(ours, theirs); err != nil {
		return reset(err)
	}
	return nil, nil
}

func (m *multi) run(ours *branchRef, theirs *branchRef) error {
	var err error

	// first merge
	s1 := &multiStep{m, ours, theirs}
	if err = s1.RunE(25); err != nil {
		return err
	}

	// second merge
	var commit string
	ours.Ref.Checkout()
	if commit, err = s1.MergeCommitTree("", merge.THEIRS); err != nil {
		return err
	}
	if m.Strategy == merge.OURS {
		if commit, err = s1.MergeCommitTree(commit, merge.OURS); err != nil {
			return err
		}
	}

	// update branch ref
	return util.RunAll(
		git.Command("update-ref", m.HeadRef.Info.RefName, commit).Run,
		git.Command("reset", "--hard", "HEAD").Run,
		git.RemoveUntracked,
		m.HeadRef.Checkout,
		ours.Ref.Delete,
		theirs.Ref.Delete,
	)
}

func (m *multiStep) RunE(findRenames int) error {
	return util.RunAll(
		NewStep(m.ours, m.theirs, findRenames).Run,
		NewStep(m.theirs, m.ours, findRenames).Run,
		m.update,
	)
}

func (m *multiStep) MergeCommitTree(head string, which merge.Strategy) (commit string, err error) {
	pick := m.theirs
	mergeHead := pick.Ref.ShortSha()
	switch which {
	case merge.OURS:
		pick = m.ours
	case merge.THEIRS:
		head = m.ours.Ref.ShortSha()
	}
	tree := pick.Ref.TreeRef()
	commitCmd := git.Command("commit-tree", tree,
		"-p", head,
		"-p", mergeHead,
		"-m", getMsg(pick.OursName))
	if commit, err = commitCmd.Output(); err == nil {
		err = git.Command("update-ref", m.ours.Ref.Info.RefName, commit).Run()
	}
	return
}

func (m *multiStep) update() error {
	return util.RunAll(
		m.ours.Update,
		m.theirs.Update,
	)
}
