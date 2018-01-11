package tree

import (
	"github.com/evovetech/got/file"
	"github.com/evovetech/got/git"
	"github.com/evovetech/got/git/object"
	"github.com/evovetech/got/git/types"
)

type Entry struct {
	git.Object

	Mode string
	Path file.Path
}

func newEntry(match []string) Entry {
	return Entry{
		Object: object.New(
			types.Id(match[3]),
			types.Parse(match[2]),
		),
		Mode: match[1],
		Path: file.GetPath(match[4]),
	}
}
