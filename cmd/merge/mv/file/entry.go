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
	File() (File, bool)
	Dir() (Dir, bool)

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

func (e *entry) Key() Path {
	return e.path.Init()
}

func (e *entry) File() (File, bool) {
	f, ok := e.value.(File)
	return f, ok
}

func (e *entry) Dir() (Dir, bool) {
	d, ok := e.value.(Dir)
	return d, ok
}

func (e *entry) setPath(path Path) {
	e.path = path
}
