package mv

import (
	"github.com/evovetech/got/cmd/merge/mv/file"
	"github.com/evovetech/got/git"
	"github.com/evovetech/got/log"
	"regexp"
)

type Map struct {
	Root    file.Dir `json:"-"`
	Renames file.Dir `json:"-"`
	Reverse file.Dir `json:"-"`
}

var reAdd = regexp.MustCompile("^A\\s+(.*)$")
var reDel = regexp.MustCompile("^D\\s+(.*)$")
var reRename = regexp.MustCompile("^R\\s+(.*)\\s+->\\s+(.*)$")
var reSrc = regexp.MustCompile("^(.*)/?src/(.*)$")

func NewMap() *Map {
	m := &Map{
		Root:    file.NewRoot(),
		Renames: file.NewRoot(),
		Reverse: file.NewRoot(),
	}
	return m
}

func (m *Map) Run() []Rename {
	for _, status := range git.Command("status", "-s").OutputLines() {
		switch {
		case reAdd.MatchString(status):
			match := reAdd.FindStringSubmatch(status)
			m.add(match[1], file.Add)
		case reDel.MatchString(status):
			match := reDel.FindStringSubmatch(status)
			m.add(match[1], file.Del)
		case reRename.MatchString(status):
			match := reRename.FindStringSubmatch(status)
			if mv, ok := file.ParseMove(match[1], match[2]); ok {
				//log.Println(mvPath.String())
				m.Renames.PutFile(mv.String(), file.Mv)
			}
		}
	}
	return m.parse()
}

func (m *Map) parse() (renames []Rename) {
	// TODO:

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

func (m *Map) add(f string, typ file.Type) {
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
	m.Root.PutFile(f, typ)
	reverse := file.GetPath(f)
	reverse.Reverse()
	m.Reverse.PutFile(reverse.String(), typ)
}
