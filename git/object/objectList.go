package object

import "github.com/evovetech/got/util"

type List []Object

func (list *List) Append(obj Object) {
	*list = append(*list, obj)
}

func (list List) String() string {
	return util.String(list)
}
