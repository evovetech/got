package file

import (
	"fmt"
	"github.com/evovetech/got/log"
)

type Entry interface {
	fmt.Stringer
	Key() Path
	Value() interface{}
	IsDir() bool
	Path() Path

	// private
	setKey(key Path)
	log(l *log.Logger)
}

type entry struct {
	key   Path
	value interface{}
}

func (e *entry) Key() Path {
	return e.key.Init()
}

func (e *entry) Value() interface{} {
	return e.value
}

func (e *entry) Path() Path {
	return e.Key()
}

func (e *entry) setKey(key Path) {
	e.key = key
}
