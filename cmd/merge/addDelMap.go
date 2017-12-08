package merge

import "strings"

type AddDelMap map[string]*MvGroup
type AddDelType int

func (t AddDelType) String() string {
	switch t {
	case Add:
		return "A"
	case Del:
		return "D"
	}
	return "?"
}

const (
	Add AddDelType = iota
	Del
)

func (m AddDelMap) Parse() (errs []*MvGroup, pairs []MvPair) {
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

func (m AddDelMap) do(match []string, typ AddDelType) FilePath {
	file := strings.Trim(match[1], "\"")
	fp := GetFilePath(file)
	fName := fp.LoName()
	mv, ok := m[fName]
	if !ok {
		mv = &MvGroup{FileName: fName}
		m[fName] = mv
	}
	switch typ {
	case Add:
		mv.Add(fp)
	case Del:
		mv.Del(fp)
	}
	return fp
}
