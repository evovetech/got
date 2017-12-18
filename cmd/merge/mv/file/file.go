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
}

type file struct {
	entry
}

func GetFile(file string, typ Type) (Path, File) {
	path := GetPath(file)
	return path.Dir(), NewFile(path.Name(), typ)
}

func NewFile(name string, typ Type) File {
	f := new(file)
	f.path = GetPath(name)
	f.value = typ
	return f
}

func (f *file) Name() string {
	return f.Path().Name()
}

func (f *file) Type() Type {
	return f.value.(Type)
}

func (f *file) IsDir() bool {
	return false
}

func (f file) String() string {
	return fmt.Sprintf("%s: '%s'", f.Type(), f.Name())
}

func (f file) MarshalJSON() ([]byte, error) {
	return json.Marshal(f.String())
}

func (f *file) log(l *log.Logger) {
	l.Println(f.String())
}
