package file

import (
	"fmt"
	"github.com/evovetech/got/log"
)

type Entry interface {
	fmt.Stringer
	Path() Path
	Value() interface{}
	IsDir() bool

	// private
	setPath(path Path)
	log(l *log.Logger)
}

type entry struct {
	path  Path
	value interface{}
}

func (e *entry) Value() interface{} {
	return e.value
}

func (e *entry) Path() Path {
	return e.path.Init()
}

func (e *entry) setPath(path Path) {
	e.path = path
}
