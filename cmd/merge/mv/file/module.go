package file

import (
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

func (m *module) Copy() Entry {
	if d, ok := m.dir.Copy().(*dir); ok {
		return &module{d}
	}
	return nil
}

func (m *module) String() string {
	return DirString(m)
}

func (m *module) log(l *log.Logger) {
	logDir(l, "module", m)
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
