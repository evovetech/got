package object

import (
	"github.com/evovetech/got/git"
)

func LsTree(sha Id) (l List) {
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
