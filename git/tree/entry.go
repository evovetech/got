package tree

import (
	"github.com/evovetech/got/file"
	"github.com/evovetech/got/git/types"
	"github.com/evovetech/got/util"
)

type EntryList []Entry

func (list *EntryList) Append(e Entry) {
	*list = append(*list, e)
}

func (list EntryList) String() string {
	return util.String(list)
}

type Entry struct {
	Mode string
	Kind types.Type
	Sha  types.Sha
	Path file.Path
}

func newEntry(match []string) Entry {
	return Entry{
		Mode: match[1],
		Kind: types.Parse(match[2]),
		Sha:  types.Sha(match[3]),
		Path: file.GetPath(match[4]),
	}
}
