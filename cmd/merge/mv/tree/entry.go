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
	Path() file.Path

	// private
	setKey(key file.Path)
	log(l *log.Logger)
}

type entry struct {
	key   file.Path
	value interface{}
}

func (e *entry) Key() file.Path {
	return e.key.Init()
}

func (e *entry) Value() interface{} {
	return e.value
}

func (e *entry) Path() file.Path {
	return e.Key()
}

func (e *entry) setKey(key file.Path) {
	e.key = key
}
