package git

import "github.com/evovetech/got/util"

type ObjectList []Object

func (list *ObjectList) Append(obj Object) {
	*list = append(*list, obj)
}

func (list ObjectList) String() string {
	return util.String(list)
}
