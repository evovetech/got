package mv

import "strings"

type AddDelMap map[string]*Group

func (m AddDelMap) parse() ([]*Group, []Rename) {
	var errs []*Group
	var pairs []Rename
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

func (m AddDelMap) do(file string, typ Type) FilePath {
	file = strings.Trim(file, "\"")
	fp := GetFilePath(file)
	fName := fp.LoName()
	mv := m.getOrCreate(fName)
	switch {
	case typ.HasFlag(Add):
		mv.Add(fp)
	case typ.HasFlag(Del):
		mv.Del(fp)
	}
	return fp
}

func (m AddDelMap) getOrCreate(fName string) *Group {
	g, ok := m[fName]
	if !ok {
		g = &Group{FileName: fName}
		m[fName] = g
	}
	return g
}
