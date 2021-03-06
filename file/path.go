/*
 * Copyright 2018 evove.tech
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package file

import (
	"encoding/json"
	"github.com/evovetech/got/util"
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
	file = removeEnds(file, '"')
	p := ospath.ToSlash(file)
	clean := spath.Clean(p)
	return Path(strings.Split(clean, "/"))
}

func removeEnds(str string, rm rune) string {
	char := string(rm)
	charLen, strLen := len(char), len(str)
	begin, end := charLen, strLen-charLen
	if end > begin &&
		str[:begin] == char &&
		str[end:] == char {
		return str[begin:end]
	}
	return str
}

func JoinPaths(paths ...Path) Path {
	switch len(paths) {
	case 0:
		return RootPath()
	case 1:
		return paths[0].Copy()
	default:
		return paths[0].Append(paths[1:]...)
	}
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

func (p Path) Mkdirs() error {
	return util.Mkdirs(p.OsPath())
}

func (p Path) Append(paths ...Path) Path {
	var m int
	var join = make(Path, p.SizeOf(paths...))
	add := func(path Path) {
		if s := path.Size(); s > 0 {
			n := m + s
			copy(join[m:n], path)
			m = n
		}
	}
	add(p)
	for _, path := range paths {
		add(path)
	}
	return join
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

func (p Path) Len() int {
	return p.Size()
}

func (p Path) Size() int {
	if p.IsRoot() {
		return 0
	}
	return len(p)
}

func (p Path) SizeOf(paths ...Path) int {
	size := p.Size()
	for _, path := range paths {
		size += path.Size()
	}
	return size
}

func (p Path) Reverse() {
	for m, n := 0, len(p)-1; m < n; m++ {
		temp := p[m]
		p[m] = p[n]
		p[n] = temp
		n--
	}
}

func (p Path) ReverseCopy() Path {
	cp := p.Copy()
	cp.Reverse()
	return cp
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
