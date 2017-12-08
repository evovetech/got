package merge

import (
	"github.com/evovetech/got/util"
	"strings"
)

type FileDir struct {
	FilePath
}

func GetFileDir(dir string) FileDir {
	return FileDir{GetFilePath(dir)}
}

func (d FileDir) Mkdirs() error {
	return util.Mkdirs(d.actual)
}

func (d FileDir) Base() string {
	return d.Name()
}

func (d FileDir) MovedFrom(other FileDir) bool {
	dStr := d.slashy
	oStr := other.slashy
	// TODO:
	return d.Base() == other.Base() ||
		strings.Contains(dStr, other.Base()) ||
		strings.Contains(oStr, d.Base())
}
