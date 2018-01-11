package object

import (
	"github.com/evovetech/got/git"
	"github.com/evovetech/got/git/types"
)

func Ls(sha types.Id) (l List) {
	git.Command("ls-tree", sha.String()).ForEachLine(func(line string) error {
		if match := reTreeLine.FindStringSubmatch(line); match != nil {
			e := newEntry(match)
			l.Append(e)
		}
		return nil
	})
	return
}
