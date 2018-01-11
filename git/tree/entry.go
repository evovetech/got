package tree

import (
	"github.com/evovetech/got/file"
	"github.com/evovetech/got/git/object"
)

type Entry struct {
	object.Object

	Mode string
	Path file.Path
}

func newEntry(match []string) Entry {
	return Entry{
		Object: object.New(
			object.Id(match[3]),
			object.Parse(match[2]),
		),
		Mode: match[1],
		Path: file.GetPath(match[4]),
	}
}
