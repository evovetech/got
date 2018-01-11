package tree

import "github.com/evovetech/got/util"

type EntryList []Entry

func (list *EntryList) Append(e Entry) {
	*list = append(*list, e)
}

func (list EntryList) String() string {
	return util.String(list)
}
