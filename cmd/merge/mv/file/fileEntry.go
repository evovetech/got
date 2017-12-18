package file

import (
	"github.com/evovetech/got/log"
)

type FileEntry interface {
	Entry
	File() File
}

type fileEntry struct {
	entry
}

func (e *fileEntry) IsDir() bool {
	return false
}

func (e *fileEntry) File() File {
	return e.value.(File)
}

func (e fileEntry) String() string {
	return e.File().String()
}

func (e *fileEntry) log(l *log.Logger) {
	l.Println(e.String())
}

func NewFileEntry(f File) FileEntry {
	e := new(fileEntry)
	e.key = GetPath(f.Name)
	e.value = f
	return e
}
