package mv

import (
	"github.com/evovetech/got/cmd/merge/mv/file"
	"github.com/evovetech/got/git"
	"github.com/evovetech/got/log"
	"regexp"
)

type Map struct {
	AddDelMap `json:"-"`
	Projects  Projects
	Root      file.Dir `json:"-"`
	Renames   file.Dir `json:"-"`
	Reverse   file.Dir `json:"-"`
}

var reAdd = regexp.MustCompile("^A\\s+(.*)$")
var reDel = regexp.MustCompile("^D\\s+(.*)$")
var reRename = regexp.MustCompile("^R\\s+(.*)\\s+->\\s+(.*)$")
var reSrc = regexp.MustCompile("^(.*)/?src/(.*)$")

func NewMap() *Map {
	m := &Map{
		AddDelMap: make(AddDelMap),
		Projects:  make(Projects),
		//Files:     make(map[string]file.Dir),
		Root:    file.NewRoot(),
		Renames: file.NewRoot(),
		Reverse: file.NewRoot(),
	}
	return m
}

func (m *Map) Run() ([]*Group, []Rename) {
	for _, status := range git.Command("status", "-s").OutputLines() {
		switch {
		case reAdd.MatchString(status):
			match := reAdd.FindStringSubmatch(status)
			m.do(match[1], Add)
		case reDel.MatchString(status):
			match := reDel.FindStringSubmatch(status)
			m.do(match[1], Del)
		case reRename.MatchString(status):
			match := reRename.FindStringSubmatch(status)
			move := file.NewMove(
				file.GetPath(match[1]),
				file.GetPath(match[2]),
			)
			if mvPath, ok := move.Parse(); ok {
				//log.Println(mvPath.String())
				m.Renames.PutFile(mvPath.String(), file.Mv)
			}
		}
	}
	return m.parse()
}

func (m *Map) parse() (groups []*Group, renames []Rename) {
	// TODO:
	for _, p := range m.Projects {
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
					//log.Printf("adding module[%s] = %s", m.Name, f)
					break
				}
			}
			if index != -1 {
				others = append(others[:index], others[index+1:]...)
			}
		}
		p.Others = others
	}

	//log.Printf("node1: %s", util.String(m.Root))
	//log.Printf("add: %s", util.String(m.Add))
	//log.Printf("add: %s", util.String(m.Add))
	//log.Printf("files: %s", util.String(m.Files))

	//for it := m.Root.DeepIterator(); it.Next(); {
	//	log.Println(it.FullPath())
	//}

	log.Println(m.Reverse.String())
	//log.Println(m.Mvs.String())
	//for _, mod := range m.Root.AllModules() {
	//	log.Print(mod.String())
	//}
	return
}

func (m *Map) getProject(dir DirPath) *Project {
	p, ok := m.Projects[dir]
	if !ok {
		p = &Project{
			Name:    dir,
			Modules: make(Modules),
		}
		m.Projects[dir] = p
	}
	return p
}

func (m *Map) add(fp FilePath, typ Type) {
	//if node := file.ParseString(fp.slashy); node != nil {
	//	m.Root.Add(node)
	//}
	//path, f := file.GetFile(fp.actual, file.Type(typ))
	//dir, ok := m.Files[f.Name]
	//if !ok {
	//	dir = file.NewDir()
	//	m.Files[f.Name] = dir
	//}
	//dir.Add(path, f)
	//m.Root.AddFile(fp.actual, file.Type(typ))
	m.Root.PutFile(fp.actual, file.Type(typ))
	reverse := file.GetPath(fp.actual)
	reverse.Reverse()
	m.Reverse.PutFile(reverse.String(), file.Type(typ))

	if src := parseSrc(fp); src != nil {
		src.Type = typ
		p := m.getProject(src.Project)
		m := p.getModule(src.Module)
		m.addSrc(*src)
		//} else if match := reProject.FindStringSubmatch(fp.slashy); match != nil {
		//	f := newProjectFile(fp, match[1])
		//	f.Type = typ
		//	f.RelPath = GetFilePath((match[2]))
		//	p := m.getProject(f.Project)
		//	p.Other = append(p.Other, *f)
	} else {
		f := newProjectFile(fp, fp.slashy)
		f.Type = typ
		//log.Printf("newProjectFile -> %s", util.String(f))
		p := m.getProject(f.Project)
		p.Others = append(p.Others, *f)
	}
}

func (m *Map) do(file string, typ Type) {
	fp := m.AddDelMap.do(file, typ)
	m.add(fp, typ)
}
