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

package git

import (
	"encoding/json"
	"regexp"
)

var (
	commits      CommitMap
	reCommitLine = regexp.MustCompile("^(\\w+)\\s+(.*)$")
)

type CommitInfo struct {
	sha     string
	tree    string
	parents map[string]bool

	children CommitMap
}

func GetCommits(commit string, num int) CommitMap {
	if ref, err := ParseRef(commit); err == nil {
		if info := GetCommitInfo(ref.Commit.Full); info != nil {
			return info.Populate(num - 1)
		}
	}
	return nil
}

func GetCommitInfo(sha string) *CommitInfo {
	info, _ := commits.getOrCreate(sha)
	return info
}

func (info *CommitInfo) Tree() string {
	return info.tree
}

func (info *CommitInfo) Parents() map[string]bool {
	return info.parents
}

func (info *CommitInfo) Children() CommitMap {
	return info.children.init()
}

func (info *CommitInfo) Populate(n int) CommitMap {
	info.populate(n)
	return info.findRoots()
}

type pTypes []*pType
type pType struct {
	info    *CommitInfo
	created bool
}

func (p *pTypes) append(types ...*pType) {
	*p = append(*p, types...)
}

func (info *CommitInfo) populate(n int) (parents pTypes) {
	if n <= 0 {
		return
	}

	for sha := range info.parents {
		p, created := commits.getOrCreate(sha)
		p.children.put(info)
		parents.append(&pType{p, created})
		parents.append(p.populate(n - 1)...)
	}
	return
}

func (info CommitInfo) MarshalJSON() ([]byte, error) {
	return json.Marshal(info.children)
}

func (info *CommitInfo) findRoots() (r CommitMap) {
	for sha := range info.parents {
		if p, found := commits.get(sha); found {
			if gr := p.findRoots(); len(gr) > 0 {
				r.putAll(gr)
			} else {
				r.put(p)
			}
		}
	}
	return
}

//
//func findRoots(commit string) {
//	c, _ := commits.getOrCreate()
//}

type CommitMap map[string]*CommitInfo

func (m *CommitMap) init() CommitMap {
	var v CommitMap
	if v = *m; v == nil {
		v = make(CommitMap)
		*m = v
	}
	return v
}

func (m *CommitMap) put(info *CommitInfo) {
	if !m.find(info.sha) {
		m.init()[info.sha] = info
	}
}

func (m *CommitMap) putAll(cm CommitMap) {
	for _, v := range cm {
		m.put(v)
	}
}

func (m *CommitMap) find(commit string) bool {
	if cm := *m; cm != nil {
		for _, info := range cm {
			if info.sha == commit || info.children.find(commit) {
				return true
			}
		}
	}
	return false
}

func (m *CommitMap) add(commit string) {
	if !m.find(commit) {
		m.init()[commit] = nil
	}
}

func (m *CommitMap) getOrCreate(commit string) (info *CommitInfo, created bool) {
	var found = false
	if info, found = m.get(commit); !found {
		info, created = newCommitInfo(commit), true
		m.init()[commit] = info
	}
	return
}

func (m *CommitMap) get(commit string) (*CommitInfo, bool) {
	if cm := *m; cm != nil {
		if info, ok := cm[commit]; ok && info != nil {
			return info, true
		}
	}
	return nil, false
}

type commitInfo struct {
	sha      string
	tree     string
	children map[string]*commitInfo
}

type CommitTree struct {
	roots map[string]*commitInfo
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

func newCommitInfo(sha string) (info *CommitInfo) {
	ci := func() *CommitInfo { return initInfo(&info, sha) }
	Command("cat-file", "-p", sha).ForEachLine(func(line string) error {
		if match := reCommitLine.FindStringSubmatch(line); match != nil {
			switch match[1] {
			case "tree":
				ci().tree = match[2]
			case "parent":
				ci().parents[match[2]] = true
			}
		}
		return nil
	})
	return
}

func initInfo(ptr **CommitInfo, sha string) (info *CommitInfo) {
	if info = *ptr; info == nil {
		info = &CommitInfo{
			sha:     sha,
			parents: make(map[string]bool),
		}
		*ptr = info
	}
	return
}
