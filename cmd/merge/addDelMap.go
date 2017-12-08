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

func (m AddDelMap) parse() (errs []*MvGroup, pairs []MvPair) {
	for _, ad := range m {
		if !ad.IsValid() {
			//log.Printf("inValid: %s", ad)
			continue
		}
		logMap := make(map[string]interface{})
		logMap["FileName"] = ad.FileName
		err, mvs := ad.parse(true)
		if len(mvs) > 0 {
			pairs = append(pairs, mvs...)
			logMap["Moves"] = mvs
		}
		if err != nil {
			errs = append(errs, err)
			logMap["Unmoved"] = err
		}
		//log.Print(util.String(logMap))
	}
	return
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
