package object

import (
	"github.com/evovetech/got/file"
	"github.com/evovetech/got/git/types"
)

type Entry struct {
	Object

	Mode string
	Path file.Path
}

func newEntry(match []string) Entry {
	return Entry{
		Object: New(
			types.Id(match[3]),
			types.Parse(match[2]),
		),
		Mode: match[1],
		Path: file.GetPath(match[4]),
	}
}
