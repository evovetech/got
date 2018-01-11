package object

import (
	"github.com/evovetech/got/file"
)

type TreeEntry struct {
	Object

	Mode string
	Path file.Path
}

func newTreeEntry(match []string) TreeEntry {
	return TreeEntry{
		Object: Parse(
			Id(match[3]),
			ParseType(match[2]),
		),
		Mode: match[1],
		Path: file.GetPath(match[4]),
	}
}
