package merge

import (
	"encoding/json"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type FilePath struct {
	actual string
	slashy string
}

func GetFilePath(file string) FilePath {
	file = filepath.FromSlash(file)
	fp := filepath.Clean(file)
	return FilePath{
		actual: fp,
		slashy: filepath.ToSlash(fp),
	}
}

func (fp FilePath) ToDir() FileDir {
	return FileDir{fp}
}

func (fp FilePath) Name() string {
	return path.Base(fp.slashy)
}

func (fp FilePath) Dir() FileDir {
	return GetFileDir(filepath.Dir(fp.actual))
}

func (fp FilePath) LoName() string {
	return strings.ToLower(fp.Name())
}

func (fp FilePath) Stat() (os.FileInfo, error) {
	return os.Stat(fp.actual)
}

func (fp FilePath) IsDir() bool {
	if info, err := fp.Stat(); err == nil {
		return info.Mode().IsDir()
	}
	return false
}

func (fp FilePath) String() string {
	return fp.actual
}

func (fp FilePath) MarshalJSON() ([]byte, error) {
	return json.Marshal(fp.String())
}
