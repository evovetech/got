package file

import (
	"encoding/json"
	"os"
	spath "path"
	ospath "path/filepath"
	"strings"
)

const (
	SRC = "src"
)

type Path []string

var (
	rootPath = GetPath("")
	srcPath  = GetPath(SRC)
)

func RootPath() Path {
	return rootPath.Copy()
}

func SrcPath() Path {
	return srcPath.Copy()
}

func GetPath(file string) Path {
	p := ospath.ToSlash(file)
	clean := spath.Clean(p)
	return Path(strings.Split(clean, "/"))
}

func JoinPaths(paths ...Path) Path {
	var length int
	for _, p := range paths {
		if !p.IsRoot() {
			length += len(p)
		}
	}

	var m int
	var join = make(Path, length)
	for _, p := range paths {
		if !p.IsRoot() {
			n := m + len(p)
			copy(join[m:n], p)
			m = n
		}
	}
	return join
}

func (p *Path) Init() Path {
	var path Path
	if path = *p; len(path) == 0 {
		path = RootPath()
		*p = path
	}
	return path
}

func (p *Path) IsRoot() bool {
	return p.Init()[0] == "."
}

func (p Path) Copy() Path {
	var l int
	if l = len(p); l == 0 {
		return RootPath()
	}

	n := make(Path, l)
	copy(n, p)
	return n
}

func (p Path) Equals(o Path) bool {
	return p.String() == o.String()
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

func (p Path) Append(paths ...Path) Path {
	paths = append([]Path{p}, paths...)
	return JoinPaths(paths...)
}

func (p Path) CopyWithParent(parent Path) Path {
	return JoinPaths(parent, p)
}

func (p Path) IndexOf(segment string) (int, bool) {
	for i, l := 0, len(p); i < l; i++ {
		if p[i] == segment {
			return i, true
		}
	}
	return -1, false
}

func (p Path) SrcIndex() (int, bool) {
	return p.IndexOf(SRC)
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
