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

package commit

import (
	"github.com/evovetech/got/collect"
	"github.com/evovetech/got/git"
	"github.com/evovetech/got/types"
	"regexp"
)

var (
	reCommitLine = regexp.MustCompile("^(\\w+)\\s+(.*)$")
)

type Info struct {
	sha     types.Sha
	tree    types.Sha
	parents collect.ShaList
}

func NewInfoFromRef(ref git.Ref) *Info {
	return NewInfo(types.Sha(ref.Commit.Full))
}

func NewInfo(sha types.Sha) (info *Info) {
	getInfo := infoGetter(&info, sha)
	git.Command("cat-file", "-p", sha.String()).ForEachLine(func(line string) error {
		if match := reCommitLine.FindStringSubmatch(line); match != nil {
			switch match[1] {
			case "tree":
				getInfo().tree = types.Sha(match[2])
			case "parent":
				getInfo().parents.Append(types.Sha(match[2]))
			}
		}
		return nil
	})
	return
}

func (info *Info) Sha() types.Sha {
	return info.sha
}

func (info *Info) Tree() types.Sha {
	return info.tree
}

func (info *Info) Parents() collect.ShaList {
	return info.parents
}

func (info *Info) FirstParent() *Info {
	if len(info.parents) == 0 {
		return nil
	}
	return NewInfo(info.parents[0])
}

type InfoGetter func() (*Info, bool)

type Group []*List
type List struct {
	info *Info
	next *List
}

func NextParentGetter(refs ...git.Ref) InfoGetter {
	var size = len(refs)
	var commits = make([]*Info, size)
	for i, ref := range refs {
		commits[i] = NewInfoFromRef(ref)
	}
	var i int
	return func() (*Info, bool) {
		return getParentInfo(&i, &commits)
	}
}

func FindForkCommit(refs ...git.Ref) (*Info, bool) {
	var commits collect.ShaCounterSet
	var nextParent = NextParentGetter(refs...)
	var target = len(refs)
	for {
		if next, ok := nextParent(); ok {
			if n := commits.Increment(next.Sha()); n >= target {
				return next, true
			}
		} else {
			return nil, false
		}
	}
}

func infoGetter(ptr **Info, sha types.Sha) func() *Info {
	return func() (info *Info) {
		if info = *ptr; info == nil {
			info = &Info{sha: sha}
			*ptr = info
		}
		return
	}
}

func getParentInfo(iPtr *int, commitsPtr *[]*Info) (*Info, bool) {
	i, commits := *iPtr, *commitsPtr
	var size int
	if size = len(commits); size == 0 {
		return nil, false
	}

	i = i % size
	if c := commits[i]; c == nil {
		*commitsPtr = append(commits[:i], commits[i+1:]...)
		return getParentInfo(iPtr, commitsPtr)
	} else {
		commits[i] = c.FirstParent()
		*iPtr = i + 1
		return c, true
	}
}
