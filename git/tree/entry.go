package tree

import (
	"github.com/evovetech/got/file"
	"github.com/evovetech/got/git/object"
)

type Entry struct {
	Mode string
	Kind object.Kind
	Sha  object.Id
	Path file.Path
}

func newEntry(match []string) Entry {
	return Entry{
		Mode: match[1],
		Kind: object.Parse(match[2]),
		Sha:  object.Id(match[3]),
		Path: file.GetPath(match[4]),
	}
}
