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
	"github.com/evovetech/got/git/object"
)

type Getter func() (Commit, bool)

func NextParentGetter(refs ...git.Ref) Getter {
	var size = len(refs)
	var commits = make([]Commit, size)
	for i, ref := range refs {
		commits[i] = New(object.Id(ref.Commit.Full))
	}
	var i int
	return func() (Commit, bool) {
		return getCommit(&i, &commits)
	}
}

func FindForkCommit(refs ...git.Ref) (Commit, bool) {
	var commits collect.ShaCounterSet
	var nextParent = NextParentGetter(refs...)
	var target = len(refs)
	for {
		if next, ok := nextParent(); ok {
			if n := commits.Increment(next.Id()); n >= target {
				return next, true
			}
		} else {
			return nil, false
		}
	}
}

func getCommit(iPtr *int, commitsPtr *[]Commit) (Commit, bool) {
	i, commits := *iPtr, *commitsPtr
	var size int
	if size = len(commits); size == 0 {
		return nil, false
	}

	i = i % size
	if c := commits[i]; c == nil {
		*commitsPtr = append(commits[:i], commits[i+1:]...)
		return getCommit(iPtr, commitsPtr)
	} else {
		commits[i] = c.Parents().Value()
		*iPtr = i + 1
		return c, true
	}
}
