package mv

import (
	"encoding/json"
	"strings"
)

type ProjectFile struct {
	File

	Project DirPath
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
	pf.Project = GetDirPath(dir)
	pf.RelPath = GetFilePath(relPath)
	return pf
}

type Project struct {
	Name    DirPath       `json:"-"`
	Modules Modules       `json:",omitempty"`
	Others  []ProjectFile `json:",omitempty"`
}

func (p *Project) getModule(dir DirPath) *Module {
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

type Projects map[DirPath]*Project

func (ps Projects) MarshalJSON() ([]byte, error) {
	projects := make(map[string]*Project)
	for k, v := range ps {
		projects[k.slashy] = v
	}
	return json.Marshal(projects)
}
