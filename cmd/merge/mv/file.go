package mv

import (
	"encoding/json"
	"fmt"
)

type File struct {
	FilePath

	Type    Type
	RelPath FilePath
}

func (f File) String() string {
	return fmt.Sprintf("%s: '%s'", f.Type, f.RelPath)
}

func (f File) MarshalJSON() ([]byte, error) {
	return json.Marshal(f.String())
}
