package merge

import (
	"encoding/json"
	"fmt"
)

type MvFile struct {
	FilePath

	Type    AddDelType
	RelPath FilePath
}

func (f MvFile) String() string {
	return fmt.Sprintf("%s: '%s'", f.Type, f.RelPath)
}

func (f MvFile) MarshalJSON() ([]byte, error) {
	return json.Marshal(f.String())
}
