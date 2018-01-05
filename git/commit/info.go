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
	tree    string
	parents collect.ShaSet
}

//func (info *CommitInfo) toTree() *CommitTree {
//	var tree CommitTree
//	tree.children = make(map[string]*CommitTree)
//
//}

/*
    ^
  ^   ^
  ^   ^
    ^
    ^
*/
/*
m1 = ShaSet
h1 = commitlist
all = map[Sha]commitlist
*/

func NewInfo(sha types.Sha) (info *Info) {
	ci := func() *Info { return initInfo(&info, sha) }
	git.Command("cat-file", "-p", sha.String()).ForEachLine(func(line string) error {
		if match := reCommitLine.FindStringSubmatch(line); match != nil {
			switch match[1] {
			case "tree":
				ci().tree = match[2]
			case "parent":
				ci().parents.Add(types.Sha(match[2]))
			}
		}
		return nil
	})
	return
}

func initInfo(ptr **Info, sha types.Sha) (info *Info) {
	if info = *ptr; info == nil {
		info = &Info{sha: sha}
		*ptr = info
	}
	return
}
