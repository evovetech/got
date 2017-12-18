package tree

import (
	"fmt"
	"github.com/evovetech/got/cmd/merge/mv/file"
	"github.com/evovetech/got/log"
)

type Value interface {
	Path() file.Path
	Value() interface{}
	IsDir() bool
}

type value struct {
	path file.Path
	val  interface{}
}

func (v *value) Path() file.Path {
	return v.path.Init()
}

func (v *value) Value() interface{} {
	return v.val
}

func (v *value) IsDir() bool {
	_, ok := v.val.(*Tree)
	return ok
}

type File struct {
	value
}

func (v *File) File() file.File {
	return v.val.(file.File)
}

func (v File) String() string {
	return v.File().String()
}

func NewFile(f file.File) *File {
	v := new(File)
	v.path = file.GetPath(f.Name)
	v.val = f
	return v
}

type Dir struct {
	value
}

func NewDir(path file.Path) *Dir {
	v := new(Dir)
	v.path = path
	v.val = New()
	return v
}

func (v *Dir) Tree() *Tree {
	return v.val.(*Tree)
}

func (v Dir) String() string {
	return fmt.Sprintf("dir: %s", v.path)
}

func (v *Dir) log(l *log.Logger) {
	v.Tree().log(l, v.Path())
}
