package object

import (
	"github.com/evovetech/got/git/types"
)

type Object interface {
	Type() types.Type
	Id() types.Sha
}

type object struct {
	typ types.Type
	id  types.Sha
}

func (o *object) Type() types.Type {
	return o.typ
}

func (o *object) Id() types.Sha {
	return o.id
}
