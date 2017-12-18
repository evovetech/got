package tree

import (
	"github.com/evovetech/got/cmd/merge/mv/file"
	"github.com/evovetech/got/log"
)

func (t *Tree) log(logger *log.Logger, path file.Path) {
	logger.Enter(path, func(l *log.Logger) {
		l.Println(t.Tree.String())
		for _, f := range t.Files() {
			l.Println(f.String())
		}
		for _, dir := range t.Dirs() {
			dir.Tree().log(l, dir.Key())
		}
	})
}
