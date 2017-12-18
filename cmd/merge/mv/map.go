package mv

import (
	"fmt"
	"github.com/evovetech/got/cmd/merge/mv/file"
	"github.com/evovetech/got/git"
	"github.com/evovetech/got/log"
	"github.com/evovetech/got/util"
	"os"
	"regexp"
)

type Map struct {
	AddDelMap `json:"-"`
	Renames   []Rename
	Projects  Projects
	Root      file.DirEntry `json:"-"`
	//Files     map[string]file.Dir
	//Mvs map[string]
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
		Root: file.NewRoot(),
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
			pair := new(Rename)
			pair.From = GetFilePath(match[1])
			pair.To = GetFilePath(match[2])
			m.Renames = append(m.Renames, *pair)
		}
	}
	for _, pair := range m.Renames {
		//dir := file.NewDir()
		//p1, f1 := file.GetFile(pair.From.actual, file.Del|file.Rn)
		//p2, f2 := file.GetFile(pair.To.actual, file.Add|file.Rn)
		//d1 := dir.Add(p1, f1)
		//d2 := dir.Add(p2, f2)
		//if d1 != d2 {
		//	if f1.Name != f2.Name {
		//		// TODO:
		//		break
		//	}
		//	// changed directories
		//	root := dir.Dirs()[0]
		//	v := struct {
		//		Root file.Path
		//		From file.Path
		//		To   file.Path
		//		File string
		//	}{
		//		Root: root.Name(),
		//		From: d1.Name(),
		//		To:   d2.Name(),
		//		File: f1.Name,
		//	}
		//	log.Printf("rename: %s", util.String(v))
		//}
		m.add(pair.From, Del|Rn)
		m.add(pair.To, Add|Rn)
	}
	return m.parse()
}

func (m *Map) parse() ([]*Group, []Rename) {
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
	log.Println(m.Root.String())
	pairs := m.Renames
	errs, p := m.AddDelMap.parse()
	if len(p) > 0 {
		pairs = append(pairs, p...)
	}
	fname := fmt.Sprintf("../mergey/out-%d.json", len(pairs))
	if f, err := os.OpenFile(fname, os.O_CREATE|os.O_WRONLY, 0644); err == nil {
		f.WriteString(util.String(m))
		f.Close()
	}
	return errs, pairs
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
