package file

import (
	"github.com/evovetech/got/log"
	"fmt"
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
	prefix := fmt.Sprintf("module<%s>", m.Path().String())
	l.Enter(prefix, func(_ *log.Logger) {
		//for t, n := range m.MvCount() {
		//	l.Printf("%s: %d\n", t.String(), n)
		//}
		for _, f := range m.AllFiles() {
			f.log(l)
		}
	})
}
