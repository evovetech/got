package file

import (
	"encoding/json"
	"os"
	spath "path"
	ospath "path/filepath"
	"strings"
)

type Path string

func GetPath(file string) Path {
	p := ospath.ToSlash(file)
	return Path(spath.Clean(p))
}

func (p Path) OsPath() string {
	return ospath.FromSlash(p.String())
}

func (p Path) Name() string {
	return spath.Base(string(p))
}

func (p Path) Dir() Path {
	return GetPath(spath.Dir(p.String()))
}

func (p Path) LoName() string {
	return strings.ToLower(p.Name())
}

func (p Path) Stat() (os.FileInfo, error) {
	return os.Stat(p.OsPath())
}

func (p Path) IsDir() bool {
	if info, err := p.Stat(); err == nil {
		return info.Mode().IsDir()
	}
	return false
}

func (p Path) String() string {
	return string(p)
}

func (p Path) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.String())
}
