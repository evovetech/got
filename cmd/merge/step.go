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
	"regexp"
	"path"
	"github.com/evovetech/got/log"
	"os"
	"strings"
)

var reAdd = regexp.MustCompile("^A\\s+(.*)")
var reDel = regexp.MustCompile("^D\\s+(.*)")

type Step struct {
	Branch      branchRef
	Target      branchRef
	Strategy    merge.Strategy
	FindRenames int
}

func NewStep(branch branchRef, target branchRef, strategy merge.Strategy) *Step {
	return &Step{
		Branch:   branch,
		Target:   target,
		Strategy: strategy,
	}
}

func (s *Step) MergeResetTheirs() error {
	if err := util.RunAll(s.checkout, s.updateBranchRef); err != nil {
		return err
	}

	targetTree := s.Target.Ref.ShortSha() + "^{tree}"
	err := git.Command("read-tree", targetTree).Run()
	if err != nil {
		return err
	}

	var headTree string
	if headTree, err = git.Command("write-tree").Output(); err != nil {
		return err
	}
	commitCmd := git.Command("commit-tree", headTree,
		"-p", s.Branch.Ref.ShortSha(),
		"-p", s.Target.Ref.ShortSha(),
		"-m", s.getMsg())
	var commitSha string
	if commitSha, err = commitCmd.Output(); err != nil {
		return err
	}
	if err = git.Command("update-ref", s.Branch.Ref.Info.RefName, commitSha).Run(); err != nil {
		return err
	}
	git.Command("add", ".").Run()
	return git.Command("reset", "--hard", "HEAD").Run()
}

type AddDelMap map[string]*AddDel
type AddDelType int

const (
	Add AddDelType = iota
	Del
)

func (m AddDelMap) Do(match []string, typ AddDelType) {
	file := match[1]
	fName := path.Base(file)
	fDir := path.Dir(file)
	ad, ok := m[fName]
	if !ok {
		ad = &AddDel{Fname: fName}
		m[fName] = ad
	}
	switch typ {
	case Add:
		ad.add(fDir)
		//log.Printf("A %s at %s", fName, fDir)
	case Del:
		ad.del(fDir)
		//log.Printf("D %s at %s", fName, fDir)
	}
}

type AddDel struct {
	Fname string
	Add   []string
	Del   []string
}

type MvPair struct {
	From string
	To   string
}

func (mv *MvPair) run() error {
	if err := mkdir(path.Dir(mv.To)); err != nil {
		return err
	}
	return git.Command("mv", mv.From, mv.To).Run()
}

func (ad *AddDel) String() string {
	return util.String(ad)
}

func (ad *AddDel) add(dir string) {
	ad.Add = append(ad.Add, dir)
}

func (ad *AddDel) del(dir string) {
	ad.Del = append(ad.Del, dir)
}

func (ad *AddDel) hasBoth() bool {
	return len(ad.Add) > 0 && len(ad.Del) > 0
}

func (ad *AddDel) isSize(size int) bool {
	return len(ad.Add) == size && len(ad.Del) == size
}

func (ad *AddDel) isSameSize() bool {
	return len(ad.Add) == len(ad.Del)
}

func (ad *AddDel) file(dir string) string {
	return path.Join(dir, ad.Fname)
}

func (ad *AddDel) first() MvPair {
	return MvPair{
		From: ad.file(ad.Del[0]),
		To:   ad.file(ad.Add[0]),
	}
}

func (ad *AddDel) parse(firstTry bool) (*AddDel, []MvPair) {
	if ad.isSize(1) {
		return nil, []MvPair{
			ad.first(),
		}
	}
	if firstTry {
		var pairs []MvPair
		var unadded []string
		var undeleted = ad.Del
		for _, addDir := range ad.Add {
			var index = -1
			var delDir string
			addBase := strings.ToLower(path.Base(addDir))
			for i, dir := range undeleted {
				delBase := strings.ToLower(path.Base(dir))
				if addBase == delBase ||
					strings.Contains(addBase, delBase) ||
					strings.Contains(delBase, addBase) {
					index = i
					delDir = dir
					break
				}
			}
			if index == -1 {
				unadded = append(unadded, addDir)
			} else {
				mv := MvPair{
					From: ad.file(delDir),
					To:   ad.file(addDir),
				}
				pairs = append(pairs, mv)
				undeleted = append(undeleted[:index], undeleted[index+1:]...)
			}
		}
		ad2 := &AddDel{
			Fname: ad.Fname,
			Add:   unadded,
			Del:   undeleted,
		}
		_, mv2 := ad2.parse(false)
		if len(mv2) > 0 {
			pairs = append(pairs, mv2...)
		}
		return ad2, pairs
	}

	return ad, nil
}

func (s *Step) Run() error {
	if err := s.merge(); err != nil {
		return err
	}

	var adm = getAdm()
	if len(adm) == 0 {
		return s.commit()
	}

	// abort merge, move files
	if err := moveAdmFiles(adm); err != nil {
		return err
	}

	// TODO: run again?
	return s.mergeCommit()
}

func getAdm() AddDelMap {
	var adm = make(AddDelMap)
	for _, status := range git.Command("status", "-s", "--untracked-files=all").OutputLines() {
		switch {
		case reAdd.MatchString(status):
			match := reAdd.FindStringSubmatch(status)
			adm.Do(match, Add)
		case reDel.MatchString(status):
			match := reDel.FindStringSubmatch(status)
			adm.Do(match, Del)
		}
	}
	var delKeys []string
	for k, v := range adm {
		if !v.hasBoth() {
			delKeys = append(delKeys, k)
		}
	}
	for _, k := range delKeys {
		delete(adm, k)
	}
	return adm
}

func moveAdmFiles(adm AddDelMap) error {
	// abort merge, move files
	var errs []*AddDel
	git.AbortMerge()
	for _, ad := range adm {
		err, mvs := ad.parse(true)
		if err != nil {
			errs = append(errs, err)
		}
		if len(mvs) > 0 {
			for _, mv := range mvs {
				mv.run()
			}
		}
	}
	log.Printf("errors: %s", util.String(errs))
	return git.Command("commit", "-m", "moving files to prepare for merge").Run()
}

func mkdir(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return os.MkdirAll(path, 0755)
	}
	return nil
}

func (s *Step) mergeCommit() error {
	if err := s.merge(); err != nil {
		return err
	}
	return s.commit()
}

func (s *Step) merge() error {
	if err := util.RunAll(s.checkout, s.updateBranchRef); err != nil {
		return err
	}

	// git merge --no-commit -X "$cmd" "$merge_commit"
	m := s.Target.Ref.Merge()
	m.NoCommit()
	m.Strategy = s.Strategy
	if s.FindRenames != 0 {
		m.FindRenames(s.FindRenames)
	}
	if err := m.Run(); err != nil {
		if err := resolve.Run(s.Strategy); err != nil {
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

func (s *Step) getMsg() string {
	head := s.Branch
	target := s.Target
	format := "merge %s into %s -- CONFLICTS -- resolving with %s changes"
	res := target
	if s.Strategy == merge.OURS {
		res = head
	}
	return fmt.Sprintf(format, target.Ref.ShortSha(), head.Ref.ShortSha(), res.OursName)
}
