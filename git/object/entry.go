package object

import (
	"github.com/evovetech/got/file"
)

type Entry struct {
	Object

	Mode string
	Path file.Path
}

func newEntry(match []string) Entry {
	return Entry{
		Object: New(
			Id(match[3]),
			ParseType(match[2]),
		),
		Mode: match[1],
		Path: file.GetPath(match[4]),
	}
}
