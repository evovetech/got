package git

import (
	"regexp"
	"encoding/json"
)

var (
	commits      CommitMap
	reCommitLine = regexp.MustCompile("^(\\w+)\\s+(.*)$")
)

type CommitInfo struct {
	sha     string
	tree    string
	parents CommitMap
}

func GetCommits(start string, num int) *CommitInfo {
	ref, ok := ParseRef(start)
	if ok != nil {
		return nil
	}
	head := ref.Commit.Full
	info := GetCommitInfo(head)
	if info != nil {
		info.Populate(num - 1)
	}
	return info
}

func GetCommitInfo(sha string) *CommitInfo {
	info, _ := commits.init().getOrCreate(sha)
	return info
}

func (info *CommitInfo) Tree() string {
	return info.tree
}

func (info *CommitInfo) Parents() CommitMap {
	return info.parents.init()
}

func (info *CommitInfo) Populate(num int) {
	info.Parents().populate(num)
}

func (info CommitInfo) MarshalJSON() ([]byte, error) {
	return json.Marshal(info.parents)
}

type CommitMap map[string]*CommitInfo

func (m *CommitMap) init() CommitMap {
	var v CommitMap
	if v = *m; v == nil {
		v = make(CommitMap)
		*m = v
	}
	return v
}

func (m *CommitMap) put(commit string, info *CommitInfo) {
	m.init()[commit] = info
}

func (m *CommitMap) add(commit string) {
	m.put(commit, nil)
}

func (m CommitMap) populate(num int) {
	if num <= 0 {
		return
	}
	var i, parents = 0, make([]string, len(m))
	for k := range m {
		parents[i] = k
		i++
	}
	for _, commit := range parents {
		if p := m[commit]; p == nil {
			p = GetCommitInfo(commit)
			p.Populate(num - 1)
			m[commit] = p
		}
	}
}

func (m CommitMap) getOrCreate(commit string) (info *CommitInfo, created bool) {
	var ok = false
	if info, ok = m[commit]; ok {
		return
	}
	info, created = newCommitInfo(commit), true
	m[commit] = info
	return
}

func newCommitInfo(sha string) (info *CommitInfo) {
	ci := func() *CommitInfo {
		if info == nil {
			info = &CommitInfo{sha: sha}
		}
		return info
	}
	Command("cat-file", "-p", sha).ForEachLine(func(line string) error {
		if match := reCommitLine.FindStringSubmatch(line); match != nil {
			switch match[1] {
			case "tree":
				ci().tree = match[2]
			case "parent":
				ci().parents.add(match[2])
			}
		}
		return nil
	})
	return
}
