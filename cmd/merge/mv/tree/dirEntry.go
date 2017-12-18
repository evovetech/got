package tree

import (
	"bytes"
	"github.com/evovetech/got/cmd/merge/mv/file"
	"github.com/evovetech/got/log"
)

type DirEntry interface {
	Entry
	Tree() *Tree
}

type dirEntry struct {
	entry
}

func NewRoot() DirEntry {
	return NewDirEntry(file.GetPath(""))
}

func NewDirEntry(path file.Path) DirEntry {
	e := new(dirEntry)
	e.key = path
	e.value = newTree()
	return e
}

func (e *dirEntry) Tree() *Tree {
	return e.value.(*Tree)
}

func (e *dirEntry) String() string {
	var buf bytes.Buffer
	l := log.NewBufLogger(&buf)
	e.log(l)
	return buf.String()
}

func (e *dirEntry) log(logger *log.Logger) {
	logger.Enter(e.Key(), func(l *log.Logger) {
		e.Tree().log(l)
	})
}
