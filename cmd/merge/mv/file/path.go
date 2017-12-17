package file

import (
	"encoding/json"
	"os"
	spath "path"
	ospath "path/filepath"
	"strings"
)

type Path []string

func GetPath(file string) Path {
	p := ospath.ToSlash(file)
	clean := spath.Clean(p)
	return Path(strings.Split(clean, "/"))
}

func (p *Path) Init() Path {
	var path Path
	if path = *p; len(path) == 0 {
		path = GetPath("")
		*p = path
	}
	return path
}

func (p Path) IsRoot() bool {
	return p[0] == "."
}

func (p Path) OsPath() string {
	return ospath.Join(p...)
}

func (p Path) SPath() string {
	return spath.Join(p...)
}

func (p Path) Name() string {
	return spath.Base(p.SPath())
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
	return p.SPath()
}

func (p Path) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.String())
}

func (p Path) IndexMatch(val Path) (int, bool) {
	var max int
	if max = len(p); max == 0 {
		return -1, false
	}
	if v := len(val); v < max {
		if v == 0 {
			return -1, false
		}
		max = v
	}
	if p[0] == "." {
		return -1, true
	}
	var index = -1
	for i, part := range p[:max] {
		if val[i] != part {
			break
		}
		index = i
	}
	return index, index != -1
}
