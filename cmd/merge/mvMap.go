package merge

import (
	"encoding/json"
	"fmt"
	"github.com/evovetech/got/log"
	"github.com/evovetech/got/util"
	"path/filepath"
	"regexp"
	"strings"
)

type MvMap struct {
	AddDelMap `json:"-"`
	Projects  Projects
}

var reSrc = regexp.MustCompile("^(.*)/?src/(.*)$")
var reProject = regexp.MustCompile("^([^/]*)/(.*)$")

func NewMvMap() *MvMap {
	return &MvMap{
		AddDelMap: make(AddDelMap),
		Projects:  make(Projects),
	}
}

func (mv *MvMap) Match(status string) bool {
	switch {
	case reAdd.MatchString(status):
		match := reAdd.FindStringSubmatch(status)
		mv.do(match, Add)
		return true
	case reDel.MatchString(status):
		match := reDel.FindStringSubmatch(status)
		mv.do(match, Del)
		return true
	default:
		return false
	}
}

func (mv *MvMap) Parse() ([]*MvGroup, []MvPair) {
	// TODO:
	for _, p := range mv.Projects {
		others := p.Others
		for _, m := range p.Modules {
			if len(others) == 0 {
				break
			}
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
	log.Println(util.String(mv))
	return mv.AddDelMap.Parse()
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
	var values []*Project
	for _, p := range ps {
		values = append(values, p)
	}
	return json.Marshal(values)
}

func (mods Modules) MarshalJSON() ([]byte, error) {
	var values []*Module
	for _, m := range mods {
		values = append(values, m)
	}
	return json.Marshal(values)
}

type Project struct {
	Name    FileDir
	Modules Modules
	Others  []ProjectFile
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
	Project FileDir `json:"-"`
	Name    FileDir
	Src     []ModuleFile
	Other   []ModuleFile

	re *regexp.Regexp
}

func (m *Module) parse(pf ProjectFile) *ModuleFile {
	fp := pf.RelPath
	if m.re == nil {
		pat := fmt.Sprintf("^%s/(.*)$", m.Name.slashy)
		m.re = regexp.MustCompile(pat)
	}
	if match := m.re.FindStringSubmatch(fp.slashy); match != nil {
		f := new(ModuleFile)
		f.ProjectFile = pf
		f.Module = m.Name
		f.RelPath = GetFilePath(filepath.FromSlash(match[1]))
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
	return fmt.Sprintf("%s: %s", f.Type, f.RelPath)
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
	pf.Project = GetFileDir(filepath.FromSlash(dir))
	pf.RelPath = GetFilePath(filepath.FromSlash(relPath))
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
		src.RelPath = GetFilePath(filepath.FromSlash(match[2]))
		return src
	}
	return nil
}

func (mv *MvMap) do(match []string, typ AddDelType) {
	fp := mv.AddDelMap.do(match, typ)
	if src := parseSrc(fp); src != nil {
		src.Type = typ
		p := mv.getProject(src.Project)
		m := p.getModule(src.Module)
		m.addSrc(*src)
		//} else if match := reProject.FindStringSubmatch(fp.slashy); match != nil {
		//	f := newProjectFile(fp, match[1])
		//	f.Type = typ
		//	f.RelPath = GetFilePath(filepath.FromSlash(match[2]))
		//	p := mv.getProject(f.Project)
		//	p.Other = append(p.Other, *f)
	} else {
		f := newProjectFile(fp, fp.slashy)
		f.Type = typ
		p := mv.getProject(f.Project)
		p.Others = append(p.Others, *f)
	}
}
