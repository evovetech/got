package mv

import (
	"encoding/json"
	"github.com/evovetech/got/util"
	"strings"
)

type DirPath struct {
	FilePath
}

func GetDirPath(dir string) DirPath {
	return DirPath{GetFilePath(dir)}
}

func (d DirPath) Mkdirs() error {
	return util.Mkdirs(d.actual)
}

func (d DirPath) Base() string {
	return d.Name()
}

func (d DirPath) MovedFrom(other DirPath) bool {
	dStr := d.slashy
	oStr := other.slashy
	// TODO:
	return d.Base() == other.Base() ||
		strings.Contains(dStr, other.Base()) ||
		strings.Contains(oStr, d.Base())
}

func (d DirPath) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.actual)
}
