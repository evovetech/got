package merge

import "strings"

type AddDelMap map[string]*MvGroup
type AddDelType uint32

func (t AddDelType) String() string {
	var str string
	if t.HasFlag(Rename) {
		str += "R"
	}
	switch {
	case t.HasFlag(Add):
		str += "A"
	case t.HasFlag(Del):
		str += "D"
	default:
		str += "?"
	}
	return str
}

func (t AddDelType) HasFlag(flag AddDelType) bool {
	return t&flag != 0
}

const (
	Add AddDelType = 1 << iota
	Del
	Rename
)

func (m AddDelMap) parse() ([]*MvGroup, []MvPair) {
	var errs []*MvGroup
	var pairs []MvPair
	for _, ad := range m {
		if !ad.IsValid() {
			continue
		}
		err, mvs := ad.parse(true)
		if err != nil {
			errs = append(errs, err)
		}
		if len(mvs) > 0 {
			pairs = append(pairs, mvs...)
		}
	}
	return errs, pairs
}

func (m AddDelMap) do(file string, typ AddDelType) FilePath {
	file = strings.Trim(file, "\"")
	fp := GetFilePath(file)
	fName := fp.LoName()
	mv, ok := m[fName]
	if !ok {
		mv = &MvGroup{FileName: fName}
		m[fName] = mv
	}
	switch {
	case typ.HasFlag(Add):
		mv.Add(fp)
	case typ.HasFlag(Del):
		mv.Del(fp)
	}
	return fp
}
