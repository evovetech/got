package tree

import (
	"bytes"
	"github.com/emirpasic/gods/trees/avltree"
	"github.com/evovetech/got/cmd/merge/mv/file"
	"github.com/evovetech/got/log"
)

type Tree struct {
	*avltree.Tree
}

func newTree() *Tree {
	return &Tree{avltree.NewWith(PathComparator)}
}

func (t *Tree) Get(path file.Path) (Entry, bool) {
	if v, ok := t.Tree.Get(path); ok {
		if val, ok := v.(Entry); ok {
			return val, true
		}
	}
	return nil, false
}

func (t *Tree) PutFile(fp string, typ file.Type) (file.Path, *Tree, FileEntry) {
	var parent *Tree
	path, f := file.GetFile(fp, typ)
	if path.IsRoot() {
		parent = t
	} else if tree := t.PutDir(path); tree != nil {
		parent = tree
	}
	if parent != nil {
		return path, parent, parent.AddFile(f)
	}
	return path, nil, nil
}

func (t *Tree) PutDir(path file.Path) *Tree {
	if path.IsRoot() {
		return t
	} else if tree, ok := t.putFloor(path); ok {
		return tree
	} else if tree, ok := t.putCeil(path); ok {
		return tree
	}
	return t.AddDir(path).Tree()
}

func (t *Tree) Add(e Entry) {
	t.Put(e.Key(), e)
}

func (t *Tree) AddFile(file file.File) FileEntry {
	e := NewFileEntry(file)
	t.Add(e)
	return e
}

func (t *Tree) AddDir(path file.Path) DirEntry {
	e := NewDirEntry(path)
	t.Add(e)
	return e
}

func (t *Tree) Files() (files []FileEntry) {
	for it := t.Iterator(); it.Next(); {
		switch e := it.Value().(type) {
		case FileEntry:
			files = append(files, e)
		}
	}
	return
}

func (t *Tree) Dirs() (dirs []DirEntry) {
	for it := t.Iterator(); it.Next(); {
		switch e := it.Value().(type) {
		case DirEntry:
			dirs = append(dirs, e)
		}
	}
	return
}

func (t *Tree) MvCount() (add int, del int) {
	for _, e := range t.Files() {
		f := e.File()
		switch {
		case f.Type.HasFlag(file.Add):
			add++
		case f.Type.HasFlag(file.Del):
			del++
		}
	}
	for _, dir := range t.Dirs() {
		a, d := dir.Tree().MvCount()
		add += a
		del += d
	}
	return
}

func (t *Tree) String() string {
	var buf bytes.Buffer
	l := log.NewBufLogger(&buf)
	t.log(l)
	return buf.String()
}

func (t *Tree) putFloor(path file.Path) (*Tree, bool) {
	if floor, ok := t.Floor(path); ok {
		return t.putNode(path, node(floor))
	}
	return nil, false
}

func (t *Tree) putCeil(path file.Path) (*Tree, bool) {
	if ceil, ok := t.Ceiling(path); ok {
		return t.putNode(path, node(ceil))
	}
	return nil, false
}

func (t *Tree) putNode(path file.Path, node *Node) (*Tree, bool) {
	if e := node.Entry(); path.Equals(e.Key()) {
		return e.(DirEntry).Tree(), true
	}
	tree, i := node.append(t, path)
	if i == -1 {
		return tree, tree != nil
	}
	dir := tree.PutDir(path[i:])
	return dir, dir != nil
}

func (t *Tree) log(l *log.Logger) {
	//l.Println(t.Tree.String())
	add, del := t.MvCount()
	if add > 0 {
		l.Printf("A: %d\n", add)
	}
	if del > 0 {
		l.Printf("D: %d\n", del)
	}
	//for _, f := range t.Files() {
	//	l.Println(f.String())
	//}
	for _, dir := range t.Dirs() {
		dir.log(l)
	}
}
