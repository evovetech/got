package file

import (
	"fmt"
	"github.com/evovetech/got/log"
)

type Module interface {
	Dir

	Src() Dir
}

type module struct {
	*dir
}

func createModule(d Dir) (Module, bool) {
	if m, ok := d.(*dir); ok {
		return &module{m}, true
	}
	return nil, false
}

func (m *module) Src() Dir {
	path := SrcPath()
	if src, ok := m.GetDir(path); ok {
		return src
	}
	return m.PutDir(path)
}

func (m *module) log(l *log.Logger) {
	prefix := fmt.Sprintf("module<%s>", m.Key().String())
	l.Enter(prefix, func(_ *log.Logger) {
		//for t, n := range m.MvCount() {
		//	l.Printf("%s: %d\n", t.String(), n)
		//}
		for _, f := range allDirFiles(m) {
			f.log(l)
		}
		for _, mod := range m.Modules() {
			mod.log(l)
		}
	})
}

//
//func (m *module) allModules() (modules []*module) {
//	for _, temp := range m.Modules() {
//		mod := temp.(*module)
//		modules = append(modules, mod)
//		for _, child := range mod.allModules() {
//			cp := *child
//			cp.setPath(cp.Path().CopyWithParent(mod.Path()))
//
//		}
//	}
//}

func allDirFiles(e Dir) (files []File) {
	files = e.Files()
	for _, d := range e.Dirs() {
		for _, f := range allDirFiles(d) {
			files = append(files, f.CopyWithParent(d.Key()))
		}
	}
	return files
}
