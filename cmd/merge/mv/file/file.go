package file

import (
	"encoding/json"
	"fmt"
)

type File struct {
	Name string
	Type Type
}

func GetFile(file string, typ Type) (Path, File) {
	path := GetPath(file)
	return path.Dir(), File{path.Name(), typ}
}

func (f File) String() string {
	return fmt.Sprintf("%s: '%s'", f.Type, f.Name)
}

func (f File) MarshalJSON() ([]byte, error) {
	return json.Marshal(f.String())
}
