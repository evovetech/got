package object

import (
	"encoding/json"
	"fmt"
	"github.com/evovetech/got/git"
	"sync"
)

type object struct {
	id   Id
	kind git.Type

	once     sync.Once
	initFunc func()
}

func New(id Id, kind git.Type) Object {
	return &object{id: id, kind: kind}
}

func (o *object) Id() Id {
	return o.id
}

func (o *object) Type() git.Type {
	return o.kind
}

func (o *object) String() string {
	return fmt.Sprintf("%s<%s>", o.kind, o.id)
}

func (o *object) MarshalJSON() ([]byte, error) {
	return json.Marshal(o.String())
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
