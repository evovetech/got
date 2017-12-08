package merge

import (
	"encoding/json"
	"fmt"
	"github.com/evovetech/got/git"
	"github.com/evovetech/got/log"
	"github.com/evovetech/got/util"
	"regexp"
	"strings"
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

	log.Println(util.String(mv.Projects))
	errs, pairs := mv.AddDelMap.parse()
	return errs, append(pairs, mv.Renames...)
}

func (mv *MvMap) getProject(dir FileDir) *Project {
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

type Projects map[FileDir]*Project
type Modules map[FileDir]*Module

func (ps Projects) MarshalJSON() ([]byte, error) {
	projects := make(map[string]*Project)
	for k, v := range ps {
		projects[k.slashy] = v
	}
	return json.Marshal(projects)
}

func (mods Modules) MarshalJSON() ([]byte, error) {
	modules := make(map[string]*Module)
	for k, v := range mods {
		modules[k.slashy] = v
	}
	return json.Marshal(modules)
}

type Project struct {
	Name    FileDir       `json:"-"`
	Modules Modules       `json:",omitempty"`
	Others  []ProjectFile `json:",omitempty"`
}

func (p *Project) getModule(dir FileDir) *Module {
	m, ok := p.Modules[dir]
	if !ok {
		m = &Module{
			Project: p.Name,
			Name:    dir,
		}
		p.Modules[dir] = m
	}
	return m
}

type Module struct {
	Project FileDir      `json:"-"`
	Name    FileDir      `json:"-"`
	Re      string       `json:",omitempty"`
	Src     []ModuleFile `json:"-"`
	Other   []ModuleFile `json:",omitempty"`

	re *regexp.Regexp
}

func (m *Module) parse(pf ProjectFile) *ModuleFile {
	fp := pf.RelPath
	if m.re == nil {
		name := m.Name.slashy
		if name == "." {
			name = ""
		} else {
			name += "/"
		}
		pat := fmt.Sprintf("^%s(.*)$", name)
		m.re = regexp.MustCompile(pat)
		m.Re = m.re.String()
	}
	if match := m.re.FindStringSubmatch(fp.slashy); match != nil {
		f := new(ModuleFile)
		f.ProjectFile = pf
		f.Module = m.Name
		f.RelPath = GetFilePath(match[1])
		return f
	}
	return nil
}

func (m *Module) addSrc(file ModuleFile) {
	m.Src = append(m.Src, file)
}

func (m *Module) addOther(file ModuleFile) {
	m.Other = append(m.Other, file)
}

type MvFile struct {
	FilePath

	Type    AddDelType
	RelPath FilePath
}

func (f MvFile) String() string {
	return fmt.Sprintf("%s: '%s'", f.Type, f.RelPath)
}

func (f MvFile) MarshalJSON() ([]byte, error) {
	return json.Marshal(f.String())
}

type ProjectFile struct {
	MvFile

	Project FileDir
}

func newProjectFile(fp FilePath, path string) *ProjectFile {
	var dir string
	var relPath = path
	if index := strings.Index(path, "/"); index != -1 {
		dir = path[:index]
		relPath = path[index+1:]
	}
	pf := new(ProjectFile)
	pf.FilePath = fp
	pf.Project = GetFileDir(dir)
	pf.RelPath = GetFilePath(relPath)
	return pf
}

type ModuleFile struct {
	ProjectFile

	Module FileDir
}

func parseSrc(fp FilePath) *ModuleFile {
	if match := reSrc.FindStringSubmatch(fp.slashy); match != nil {
		src := new(ModuleFile)
		src.ProjectFile = *newProjectFile(fp, match[1])
		src.Module = src.RelPath.ToDir()
		src.RelPath = GetFilePath(match[2])
		return src
	}
	return nil
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
		p := mv.getProject(f.Project)
		p.Others = append(p.Others, *f)
	}
}

func (mv *MvMap) do(file string, typ AddDelType) {
	fp := mv.AddDelMap.do(file, typ)
	mv.add(fp, typ)
}
