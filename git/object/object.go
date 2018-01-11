package object

import (
	"fmt"
	"github.com/evovetech/got/types"
	"sync"
)

type (
	Id   = types.Sha
	Kind = Type
)

type Object struct {
	id   Id
	kind Kind

	once     sync.Once
	initFunc func()
}

func New(id Id, kind Kind) *Object {
	return &Object{id: id, kind: kind}
}

func NewTree(id Id) *Object {
	return New(id, Tree)
}

func NewCommit(id Id) *Object {
	return New(id, Commit)
}

func NewBlob(id Id) *Object {
	return New(id, Blob)
}

func (o *Object) Id() Id {
	return o.id
}

func (o *Object) Kind() Kind {
	return o.kind
}

func (o *Object) String() string {
	return fmt.Sprintf("%s<%s>", o.kind, o.id)
}

func (o *Object) SetInitFunc(f func()) {
	if o.initFunc == nil {
		o.initFunc = f
	}
}

func (o *Object) Init() {
	o.once.Do(o.getInitFunc())
}

func (o *Object) getInitFunc() func() {
	if f := o.initFunc; f != nil {
		return f
	}
	return func() {}
}
