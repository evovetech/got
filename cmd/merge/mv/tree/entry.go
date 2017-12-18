package tree

import (
	"fmt"
	"github.com/evovetech/got/cmd/merge/mv/file"
)

type Entry interface {
	Key() file.Path
	Value() interface{}
	IsDir() bool
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

type File struct {
	entry
}

func (v *File) File() file.File {
	return v.value.(file.File)
}

func (v File) String() string {
	return v.File().String()
}

func NewFile(f file.File) *File {
	v := new(File)
	v.key = file.GetPath(f.Name)
	v.value = f
	return v
}

type Dir struct {
	entry
}

func NewDir(path file.Path) *Dir {
	v := new(Dir)
	v.key = path
	v.value = New()
	return v
}

func (v *Dir) Tree() *Tree {
	return v.value.(*Tree)
}

func (v Dir) String() string {
	return fmt.Sprintf("dir: %s", v.key)
}
