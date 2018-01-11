package object

import (
	"github.com/evovetech/got/git"
)

func lsTree(sha Id) (l List) {
	git.Command("ls-tree", sha.String()).ForEachLine(func(line string) error {
		if match := reTreeLine.FindStringSubmatch(line); match != nil {
			e := newTreeEntry(match)
			l.Append(e)
		}
		return nil
	})
	return
}

func catBlob(sha Id) ([]byte, error) {
	cmd := git.Command("cat-file", "blob", sha.String())
	return cmd.OutputBytes()
}

func catCommit(sha Id) (tree Tree, parents CommitList, err error) {
	err = git.Command("cat-file", "-p", sha.String()).ForEachLine(func(line string) error {
		if match := reCommitLine.FindStringSubmatch(line); match != nil {
			switch match[1] {
			case "tree":
				tree = NewTree(Id(match[2]))
			case "parent":
				p := NewCommit(Id(match[2]))
				parents.Append(p)
			}
		}
		return nil
	})
	return
}
