package tree

import (
	"fmt"
	"github.com/evovetech/got/cmd/merge/mv/file"
	"github.com/evovetech/got/log"
)

type Entry interface {
	fmt.Stringer
	Key() file.Path
	Value() interface{}
	IsDir() bool

	// private
	setKey(key file.Path)
	log(l *log.Logger)
}

type entry struct {
	key   file.Path
	value interface{}
}

func (v *entry) Key() file.Path {
	return v.key.Init()
}

func (v *entry) Value() interface{} {
	return v.value
}

func (v *entry) IsDir() bool {
	_, ok := v.value.(*Tree)
	return ok
}

func (v *entry) setKey(key file.Path) {
	v.key = key
}
