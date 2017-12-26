package mv

import (
	"github.com/evovetech/got/cmd/merge/mv/file"
	"github.com/evovetech/got/git"
	"github.com/evovetech/got/log"
	"regexp"
)

type Map struct {
	Root      file.Dir `json:"-"`
	Rename1   file.Dir `json:"-"`
	Rename2   file.Dir `json:"-"`
	Reverse   file.Dir `json:"-"`
	Possibles file.Dir `json:"-"`
}

var reAdd = regexp.MustCompile("^A\\s+(.*)$")
var reDel = regexp.MustCompile("^D\\s+(.*)$")
var reRename = regexp.MustCompile("^R\\s+(.*)\\s+->\\s+(.*)$")
var reSrc = regexp.MustCompile("^(.*)/?src/(.*)$")

func NewMap() *Map {
	m := &Map{
		Root:      file.NewRoot(),
		Rename1:   file.NewRoot(),
		Rename2:   file.NewRoot(),
		Reverse:   file.NewRoot(),
		Possibles: file.NewRoot(),
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
			m.add(match[1], file.Del|file.Rn)
			m.add(match[2], file.Add|file.Rn)
			if mv, ok := file.ParseMove(match[1], match[2]); ok {
				//log.Println(mvPath.String())
				m.Rename1.PutFile(mv.String(), file.Mv)
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

	//log.Println(m.Reverse.String())
	//log.Println(m.Mvs.String())
	//for _, mod := range m.Root.AllModules() {
	//	log.Print(mod.String())
	//}
	//l := log.Std
	m.addReverse()
	log.Println(m.Root)
	log.Println(m.Rename2)
	//log.Println(m.Rename1)
	log.Println(m.Possibles)

	return
}

func (m *Map) addReverse() {
	m.addReverseDir(nil, m.Reverse)
}

func (m *Map) addReverseDir(p file.Path, d file.Dir) {
	parent := file.JoinPaths(p, d.Path())
	for it := d.Iterator(); it.Next(); {
		switch e := it.Entry().(type) {
		case file.Dir:
			count := e.MvCount()
			ac := count.AllAdd()
			dc := count.AllDel()
			if ac == 0 || dc == 0 {
				addAllReverse(parent, e, m.Root)
			} else if ac == 1 && dc == 1 {
				m.addReverseRename(parent, e)
			} else {
				m.addReverseDir(parent, e)
			}
		case file.File:
			p := file.JoinPaths(parent, e.Path())
			m.Root.PutFile(p.String(), e.Type())
		}
	}
}

func (m *Map) addReverseRename(parent file.Path, d file.Dir) {
	var from, to file.Path
	for it := d.DeepIterator(); it.Next(); {
		switch e := it.Entry().(type) {
		case file.File:
			fp := file.JoinPaths(parent, it.FullPath())
			fp.Reverse()
			if e.Type().HasFlag(file.Add) {
				to = fp
			} else if e.Type().HasFlag(file.Del) {
				from = fp
			}
		}
	}
	if mv, ok := file.NewMove(from, to).Parse(); ok {
		str := mv.String()
		log.Printf("add rename: %s", str)
		m.Rename2.PutFile(str, file.Mv)
	} else {
		log.Printf("error add rename from='%s', to='%s'", from, to)
	}
}

func addAllReverse(parent file.Path, from file.Dir, to file.Dir) {
	for it := from.DeepIterator(); it.Next(); {
		switch e := it.Entry().(type) {
		case file.File:
			fp := file.JoinPaths(parent, it.FullPath())
			fp.Reverse()
			to.PutFile(fp.String(), e.Type())
		}
	}
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
	//m.Root.PutFile(f, typ)
	//log.Printf("%s: %s", typ, f)
	reverse := file.GetPath(f)
	reverse.Reverse()
	//log.Printf("%s: %s", typ, reverse.String())
	m.Reverse.PutFile(reverse.String(), typ)
}
