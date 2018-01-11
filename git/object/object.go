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

type Object interface {
	Id() Id
	Kind() Kind
	String() string

	SetInitFunc(func())
	Init()
}

type object struct {
	id   Id
	kind Kind

	once     sync.Once
	initFunc func()
}

func New(id Id, kind Kind) Object {
	return &object{id: id, kind: kind}
}

func NewTree(id Id) Object {
	return New(id, Tree)
}

func NewCommit(id Id) Object {
	return New(id, Commit)
}

func NewBlob(id Id) Object {
	return New(id, Blob)
}

func (o *object) Id() Id {
	return o.id
}

func (o *object) Kind() Kind {
	return o.kind
}

func (o *object) String() string {
	return fmt.Sprintf("%s<%s>", o.kind, o.id)
}

func (o *object) SetInitFunc(f func()) {
	if o.initFunc == nil {
		o.initFunc = f
	}
}

func (o *object) Init() {
	o.once.Do(o.getInitFunc())
}

func (o *object) getInitFunc() func() {
	if f := o.initFunc; f != nil {
		return f
	}
	return func() {}
}
