package merge

import (
	"path/filepath"
	"strings"
)

type FilePath struct {
	Path string
	Name string
	Dir  FileDir
}

func NewFilePath(path string) FilePath {
	return FilePath{
		Path: filepath.Clean(path),
		Name: strings.ToLower(filepath.Base(path)),
		Dir:  FileDir(filepath.Dir(path)),
	}
}

func (fp FilePath) String() string {
	return fp.Path
}
