package file

import (
	"bytes"
	"github.com/emirpasic/gods/trees/avltree"
	"github.com/evovetech/got/log"
)

type DirEntry interface {
	Entry

	PutFile(fp string, typ Type) (DirEntry, FileEntry)
	PutDir(path Path) DirEntry

	Files() []FileEntry
	Dirs() []DirEntry

	MvCount() (add int, del int)

	// private
	tree() *avltree.Tree
	add(e Entry)
	addDir(path Path) DirEntry
	addFile(file File) FileEntry
}

type dirEntry struct {
	entry
}

func NewRoot() DirEntry {
	return NewDirEntry(GetPath(""))
}

func NewDirEntry(path Path) DirEntry {
	e := new(dirEntry)
	e.key = path
	e.value = avltree.NewWith(PathComparator)
	return e
}

func (d *dirEntry) IsDir() bool {
	return true
}

func (d *dirEntry) String() string {
	var buf bytes.Buffer
	l := log.NewBufLogger(&buf)
	d.log(l)
	return buf.String()
}

func (d *dirEntry) tree() *avltree.Tree {
	return d.value.(*avltree.Tree)
}

func (d *dirEntry) log(logger *log.Logger) {
	logger.Enter(d.Key(), func(l *log.Logger) {
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
