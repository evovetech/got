package file

import (
	"bytes"
	"github.com/emirpasic/gods/trees/avltree"
	"github.com/evovetech/got/log"
	"fmt"
)

type Dir interface {
	Entry

	Get(path Path) (Entry, bool)
	GetDir(path Path) (Dir, bool)
	Find(path Path) (Entry, bool)
	FindDir(path Path) (Dir, bool)
	PutFile(fp string, typ Type) (Dir, File)
	PutDir(path Path) Dir

	Files() []File
	Dirs() []Dir
	Modules() []Module
	AllFiles() []File

	MvCount() TypeCount

	// private
	tree() *avltree.Tree
	insertDir(path Path) Dir
	put(e Entry)
	putDir(path Path) Dir
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

func (d *dir) log(l *log.Logger) {
	prefix := fmt.Sprintf("dir<%s>", d.Path().String())
	l.Enter(prefix, func(_ *log.Logger) {
		//l.Println(d.tree.String())
		//for t, n := range d.MvCount() {
		//	l.Printf("%s: %d\n", t.String(), n)
		//}
		//for _, f := range d.Files() {
		//	l.Println(f.String())
		//}
		for _, dir := range d.Dirs() {
			dir.log(l)
		}
	})
}
