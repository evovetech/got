package merge

import (
	"fmt"
	"github.com/evovetech/got/git"
	"github.com/evovetech/got/log"
	"github.com/evovetech/got/util"
	"os"
	"regexp"
)

type MvMap struct {
	AddDelMap `json:"-"`
	Renames   []MvPair
	Projects  Projects
}

var reSrc = regexp.MustCompile("^(.*)/?src/(.*)$")

func NewMvMap() *MvMap {
	return &MvMap{
		AddDelMap: make(AddDelMap),
		Projects:  make(Projects),
	}
}

func (mv *MvMap) Run() ([]*MvGroup, []MvPair) {
	for _, status := range git.Command("status", "-s").OutputLines() {
		switch {
		case reAdd.MatchString(status):
			match := reAdd.FindStringSubmatch(status)
			mv.do(match[1], Add)
		case reDel.MatchString(status):
			match := reDel.FindStringSubmatch(status)
			mv.do(match[1], Del)
		case reRename.MatchString(status):
			match := reRename.FindStringSubmatch(status)
			pair := new(MvPair)
			pair.From = GetFilePath(match[1])
			pair.To = GetFilePath(match[2])
			mv.Renames = append(mv.Renames, *pair)
		}
	}
	for _, pair := range mv.Renames {
		mv.add(pair.From, Del|Rename)
		mv.add(pair.To, Add|Rename)
	}
	return mv.parse()
}

func (mv *MvMap) parse() ([]*MvGroup, []MvPair) {
	// TODO:
	for _, p := range mv.Projects {
		if len(p.Others) == 0 {
			continue
		}
		others := append([]ProjectFile{}, p.Others...)
		for _, m := range p.Modules {
			var index = -1
			for i, pf := range others {
				if mf := m.parse(pf); mf != nil {
					index = i
					f := *mf
					m.addOther(f)
					log.Printf("adding module[%s] = %s", m.Name, f)
					break
				}
			}
			if index != -1 {
				others = append(others[:index], others[index+1:]...)
			}
		}
		p.Others = others
	}

	pairs := mv.Renames
	errs, p := mv.AddDelMap.parse()
	if len(p) > 0 {
		pairs = append(pairs, p...)
	}
	fname := fmt.Sprintf("../mergey/out-%d.json", len(pairs))
	if f, err := os.OpenFile(fname, os.O_CREATE|os.O_WRONLY, 0644); err == nil {
		f.WriteString(util.String(mv))
		f.Close()
	}
	return errs, pairs
}

func (mv *MvMap) getProject(dir DirPath) *Project {
	p, ok := mv.Projects[dir]
	if !ok {
		p = &Project{
			Name:    dir,
			Modules: make(Modules),
		}
		mv.Projects[dir] = p
	}
	return p
}

func (mv *MvMap) add(fp FilePath, typ AddDelType) {
	if src := parseSrc(fp); src != nil {
		src.Type = typ
		p := mv.getProject(src.Project)
		m := p.getModule(src.Module)
		m.addSrc(*src)
		//} else if match := reProject.FindStringSubmatch(fp.slashy); match != nil {
		//	f := newProjectFile(fp, match[1])
		//	f.Type = typ
		//	f.RelPath = GetFilePath((match[2]))
		//	p := mv.getProject(f.Project)
		//	p.Other = append(p.Other, *f)
	} else {
		f := newProjectFile(fp, fp.slashy)
		f.Type = typ
		log.Printf("newProjectFile -> %s", util.String(f))
		p := mv.getProject(f.Project)
		p.Others = append(p.Others, *f)
	}
}

func (mv *MvMap) do(file string, typ AddDelType) {
	fp := mv.AddDelMap.do(file, typ)
	mv.add(fp, typ)
}
