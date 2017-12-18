package tree

import (
	"bytes"
	"github.com/emirpasic/gods/trees/avltree"
	"github.com/evovetech/got/cmd/merge/mv/file"
	"github.com/evovetech/got/log"
)

type DirEntry interface {
	Entry

	PutFile(fp string, typ file.Type) (DirEntry, FileEntry)
	PutDir(path file.Path) DirEntry

	Files() []FileEntry
	Dirs() []DirEntry

	MvCount() (add int, del int)

	// private
	tree() *avltree.Tree
	add(e Entry)
	addDir(path file.Path) DirEntry
	addFile(file file.File) FileEntry
}

type dirEntry struct {
	entry
}

func NewRoot() DirEntry {
	return NewDirEntry(file.GetPath(""))
}

func NewDirEntry(path file.Path) DirEntry {
	e := new(dirEntry)
	e.key = path
	e.value = avltree.NewWith(PathComparator)
	return e
}

func (de *dirEntry) IsDir() bool {
	return true
}

func (de *dirEntry) String() string {
	var buf bytes.Buffer
	l := log.NewBufLogger(&buf)
	de.log(l)
	return buf.String()
}

func (de *dirEntry) tree() *avltree.Tree {
	return de.value.(*avltree.Tree)
}

func (de *dirEntry) log(logger *log.Logger) {
	logger.Enter(de.Key(), func(l *log.Logger) {
		//l.Println(de.tree.String())
		add, del := de.MvCount()
		if add > 0 {
			l.Printf("A: %d\n", add)
		}
		if del > 0 {
			l.Printf("D: %d\n", del)
		}
		//for _, f := range de.Files() {
		//	l.Println(f.String())
		//}
		for _, dir := range de.Dirs() {
			dir.log(l)
		}
	})
}
