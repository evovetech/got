package tree

import (
	"bytes"
	"github.com/emirpasic/gods/trees/avltree"
	"github.com/emirpasic/gods/utils"
	"github.com/evovetech/got/cmd/merge/mv/file"
	"github.com/evovetech/got/log"
)

type Tree struct {
	*avltree.Tree
}

func minInt(a, b int) int {
	if a <= b {
		return a
	}
	return b
}

func PathComparator(a, b interface{}) int {
	ap := a.(file.Path)
	bp := b.(file.Path)
	min := minInt(len(ap), len(bp))
	comp := utils.StringComparator
	for i := 0; i < min; i++ {
		if diff := comp(ap[i], bp[i]); diff != 0 {
			return diff
		}
	}
	return len(ap) - len(bp)
}

func New() *Tree {
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

func (t *Tree) PutFilePath(fp string, typ file.Type) (file.Path, *Tree, *File) {
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
	return t.AddDir(path)
}

func (t *Tree) Add(v Entry) {
	t.Put(v.Key(), v)
}

func (t *Tree) AddFile(file file.File) *File {
	v := NewFile(file)
	t.Add(v)
	return v
}

func (t *Tree) AddDir(path file.Path) *Tree {
	v := NewDir(path)
	t.Add(v)
	return v.Tree()
}

func (t *Tree) Files() (files []file.File) {
	for it := t.Iterator(); it.Next(); {
		switch v := it.Value().(type) {
		case *File:
			files = append(files, v.File())
		}
	}
	return
}

func (t *Tree) Dirs() (dirs []*Dir) {
	for it := t.Iterator(); it.Next(); {
		switch v := it.Value().(type) {
		case *Dir:
			dirs = append(dirs, v)
		}
	}
	return
}

func (t *Tree) MvCount() (add int, del int) {
	for _, f := range t.Files() {
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
	t.log(l, file.GetPath(""))
	return buf.String()
}

func (t *Tree) log(logger *log.Logger, path file.Path) {
	logger.Enter(path, func(l *log.Logger) {
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
			dir.Tree().log(l, dir.Key())
		}
	})
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
		return e.(*Dir).Tree(), true
	}
	tree, i := node.append(t, path)
	if i == -1 {
		return tree, tree != nil
	}
	dir := tree.PutDir(path[i:])
	return dir, dir != nil
}
