package merge

import (
	"github.com/evovetech/got/util"
	"strings"
)

type FileDir FilePath

func GetFileDir(dir string) FileDir {
	fp := GetFilePath(dir)
	return (FileDir)(fp)
}

func (d FileDir) ToFilePath() FilePath {
	return (FilePath)(d)
}

func (d FileDir) Mkdirs() error {
	return util.Mkdirs(d.actual)
}

func (d FileDir) Base() string {
	return d.ToFilePath().Name()
}

func (d FileDir) String() string {
	return d.ToFilePath().String()
}

func (d FileDir) MovedFrom(other FileDir) bool {
	dStr := d.slashy
	oStr := other.slashy
	// TODO:
	return d.Base() == other.Base() ||
		strings.Contains(dStr, other.Base()) ||
		strings.Contains(oStr, d.Base())
}
