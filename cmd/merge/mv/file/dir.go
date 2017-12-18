package file

import (
	"bytes"
	"github.com/emirpasic/gods/trees/avltree"
	"github.com/evovetech/got/log"
)

type Dir interface {
	Entry

	PutFile(fp string, typ Type) (Dir, File)
	PutDir(path Path) Dir

	Files() []File
	Dirs() []Dir

	MvCount() (add int, del int)

	// private
	tree() *avltree.Tree
	add(e Entry)
	addDir(path Path) Dir
}

type dir struct {
	entry
}

func NewRoot() Dir {
	return NewDir(GetPath(""))
}

func NewDir(path Path) Dir {
	e := new(dir)
	e.path = path
	e.value = avltree.NewWith(PathComparator)
	return e
}

func (d *dir) IsDir() bool {
	return true
}

func (d *dir) String() string {
	var buf bytes.Buffer
	l := log.NewBufLogger(&buf)
	d.log(l)
	return buf.String()
}

func (d *dir) tree() *avltree.Tree {
	return d.value.(*avltree.Tree)
}

func (d *dir) log(logger *log.Logger) {
	logger.Enter(d.Path(), func(l *log.Logger) {
		//l.Println(d.tree.String())
		add, del := d.MvCount()
		if add > 0 {
			l.Printf("A: %d\n", add)
		}
		if del > 0 {
			l.Printf("D: %d\n", del)
		}
		//for _, f := range d.Files() {
		//	l.Println(f.String())
		//}
		for _, dir := range d.Dirs() {
			dir.log(l)
		}
	})
}