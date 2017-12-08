package merge

import "path"

type AddDelMap map[string]*AddDel
type AddDelType int

const (
	Add AddDelType = iota
	Del
)

func (m AddDelMap) Match(status string) bool {
	switch {
	case reAdd.MatchString(status):
		match := reAdd.FindStringSubmatch(status)
		m.do(match, Add)
		return true
	case reDel.MatchString(status):
		match := reDel.FindStringSubmatch(status)
		m.do(match, Del)
		return true
	default:
		return false
	}
}

func (m AddDelMap) Parse() (errs []*AddDel, pairs []MvPair) {
	for _, ad := range m {
		if !ad.hasBoth() {
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
	return
}

func (m AddDelMap) do(match []string, typ AddDelType) {
	file := match[1]
	fName := path.Base(file)
	fDir := path.Dir(file)
	ad, ok := m[fName]
	if !ok {
		ad = &AddDel{Fname: fName}
		m[fName] = ad
	}
	switch typ {
	case Add:
		ad.add(fDir)
	case Del:
		ad.del(fDir)
	}
}
