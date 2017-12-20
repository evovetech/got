package file

import (
	"encoding/json"
	"fmt"
	"github.com/evovetech/got/log"
)

type File interface {
	Entry
	Name() string
	Type() Type
	Copy() File
	CopyWithParent(parent Path) File
}

type file struct {
	entry
}

func GetFile(file string, typ Type) (Path, File) {
	path := GetPath(file)
	return path.Dir(), NewFile(path.Name(), typ)
}

func NewFile(file string, typ Type) File {
	return NewFileWithPath(GetPath(file), typ)
}

func NewFileWithPath(path Path, typ Type) File {
	f := new(file)
	f.path = path
	f.value = typ
	return f
}

func (f *file) Name() string {
	return f.Key().Name()
}

func (f *file) Type() Type {
	return f.value.(Type)
}

func (f *file) IsDir() bool {
	return false
}

func (f *file) Copy() File {
	return NewFileWithPath(f.path.Copy(), f.Type())
}

func (f *file) CopyWithParent(parent Path) File {
	path := f.Key().CopyWithParent(parent)
	return NewFileWithPath(path, f.Type())
}

func (f file) String() string {
	return fmt.Sprintf("%s: '%s'", f.Type(), f.Key().String())
}

func (f file) MarshalJSON() ([]byte, error) {
	return json.Marshal(f.String())
}

func (f *file) log(l *log.Logger) {
	l.Println(f.String())
}
