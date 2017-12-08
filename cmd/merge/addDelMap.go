package merge

type AddDelMap map[string]*MvGroup
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

func (m AddDelMap) Parse() (errs []*MvGroup, pairs []MvPair) {
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
	return
}

func (m AddDelMap) do(match []string, typ AddDelType) {
	fp := NewFilePath(match[1])
	mv, ok := m[fp.Name]
	if !ok {
		mv = &MvGroup{FileName: fp.Name}
		m[fp.Name] = mv
	}
	switch typ {
	case Add:
		mv.Add(fp)
	case Del:
		mv.Del(fp)
	}
}
