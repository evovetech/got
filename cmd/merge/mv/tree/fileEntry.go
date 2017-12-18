package tree

import (
	"github.com/evovetech/got/cmd/merge/mv/file"
	"github.com/evovetech/got/log"
)

type FileEntry interface {
	Entry
	File() file.File
}

type fileEntry struct {
	entry
}

func (e *fileEntry) IsDir() bool {
	return false
}

func (e *fileEntry) File() file.File {
	return e.value.(file.File)
}

func (e fileEntry) String() string {
	return e.File().String()
}

func (e *fileEntry) log(l *log.Logger) {
	l.Println(e.String())
}

func NewFileEntry(f file.File) FileEntry {
	e := new(fileEntry)
	e.key = file.GetPath(f.Name)
	e.value = f
	return e
}
