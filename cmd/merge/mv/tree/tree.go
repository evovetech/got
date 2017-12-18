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

func PathComparator(a, b interface{}) int {
	return utils.StringComparator(a.(file.Path).String(), b.(file.Path).String())
}

func New() *Tree {
	return &Tree{avltree.NewWith(PathComparator)}
}

func (t *Tree) String() string {
	var buf bytes.Buffer
	l := log.NewBufLogger(&buf)
	t.log(l, file.GetPath(""))
	return buf.String()
}

func (t *Tree) Files() (files []file.File) {
	for it := t.Tree.Iterator(); it.Next(); {
		switch v := it.Value().(type) {
		case *File:
			files = append(files, v.File())
		}
	}
	return
}

func (t *Tree) Dirs() (dirs []*Dir) {
	for it := t.Tree.Iterator(); it.Next(); {
		switch v := it.Value().(type) {
		case *Dir:
			dirs = append(dirs, v)
		}
	}
	return
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
	var base *Tree
	path, f := file.GetFile(fp, typ)
	if path.IsRoot() {
		base = t
	} else if tree := t.PutDir(path); tree != nil {
		base = tree
	}
	if base != nil {
		return path, base, base.AddFile(f)
	}
	return path, nil, nil
}

func (t *Tree) PutDir(path file.Path) *Tree {
	if path.IsRoot() {
		return t
	} else if cur, ok := t.Get(path); ok {
		return cur.(*Dir).Tree()
	} else if tree, ok := t.putCeil(path); ok {
		return tree
	} else if tree, ok := t.putFloor(path); ok {
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
	tree, i := node.append(t, path)
	if i == -1 {
		return tree, tree != nil
	}
	dir := tree.PutDir(path[i:])
	return dir, dir != nil
}
