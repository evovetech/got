package merge

import (
	"regexp"
	"github.com/evovetech/got/log"
	"github.com/evovetech/got/util"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"
)

type MvMap struct {
	AddDelMap `json:"-"`
	Projects Projects
	Other    []MvFile
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
	Other   []ProjectFile
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
	Src     []SrcFile
}

func (m *Module) add(file SrcFile) {
	m.Src = append(m.Src, file)
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

type SrcFile struct {
	ProjectFile

	Module FileDir
}

func parseSrc(fp FilePath) *SrcFile {
	if match := reSrc.FindStringSubmatch(fp.slashy); match != nil {
		src := new(SrcFile)
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
		m.add(*src)
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
		p.Other = append(p.Other, *f)
	}
}
