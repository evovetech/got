package merge

import (
	"github.com/evovetech/got/util"
	"path"
	"strings"
)

type AddDelMap map[string]*AddDel
type AddDelType int

const (
	Add AddDelType = iota
	Del
)

func (m AddDelMap) Do(match []string, typ AddDelType) {
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
		//log.Printf("A %s at %s", fName, fDir)
	case Del:
		ad.del(fDir)
		//log.Printf("D %s at %s", fName, fDir)
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

type AddDel struct {
	Fname string
	Add   []string
	Del   []string
}

func (ad *AddDel) String() string {
	return util.String(ad)
}

func (ad *AddDel) add(dir string) {
	ad.Add = append(ad.Add, dir)
}

func (ad *AddDel) del(dir string) {
	ad.Del = append(ad.Del, dir)
}

func (ad *AddDel) hasBoth() bool {
	return len(ad.Add) > 0 && len(ad.Del) > 0
}

func (ad *AddDel) isSize(size int) bool {
	return len(ad.Add) == size && len(ad.Del) == size
}

func (ad *AddDel) isSameSize() bool {
	return len(ad.Add) == len(ad.Del)
}

func (ad *AddDel) file(dir string) string {
	return path.Join(dir, ad.Fname)
}

func (ad *AddDel) first() MvPair {
	return MvPair{
		From: ad.file(ad.Del[0]),
		To:   ad.file(ad.Add[0]),
	}
}

func (ad *AddDel) parse(firstTry bool) (*AddDel, []MvPair) {
	if ad.isSize(1) {
		return nil, []MvPair{
			ad.first(),
		}
	}
	if firstTry {
		var pairs []MvPair
		var unadded []string
		var undeleted = ad.Del
		for _, addDir := range ad.Add {
			var index = -1
			var delDir string
			addBase := strings.ToLower(path.Base(addDir))
			for i, dir := range undeleted {
				delBase := strings.ToLower(path.Base(dir))
				if addBase == delBase ||
					strings.Contains(addBase, delBase) ||
					strings.Contains(delBase, addBase) {
					index = i
					delDir = dir
					break
				}
			}
			if index == -1 {
				unadded = append(unadded, addDir)
			} else {
				mv := MvPair{
					From: ad.file(delDir),
					To:   ad.file(addDir),
				}
				pairs = append(pairs, mv)
				undeleted = append(undeleted[:index], undeleted[index+1:]...)
			}
		}
		var err *AddDel
		if len(unadded) > 0 {
			var mv2 []MvPair
			err, mv2 = (&AddDel{
				Fname: ad.Fname,
				Add:   unadded,
				Del:   undeleted,
			}).parse(false)
			if len(mv2) > 0 {
				pairs = append(pairs, mv2...)
			}
		}
		return err, pairs
	}

	return ad, nil
}
