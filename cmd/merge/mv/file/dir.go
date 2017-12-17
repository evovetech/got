package file

import (
	"bytes"
	"container/list"
	"encoding/json"
	"fmt"
	"github.com/emirpasic/gods/trees/avltree"
	"github.com/emirpasic/gods/utils"
	"github.com/evovetech/got/log"
)

type Dir interface {
	Path() Path
	Parent() Path
	Base() Path
	Dirs() []Dir
	Files() []File
	AddFile(string, Type) Dir
	Add(Path, File) Dir

	addPath(Path) Dir
	appendFile(File)
	//replaceDir(cur Dir, with Dir) bool
}

type DirTree struct {
	*avltree.Tree
}

func (v *DirValue) log(l *log.Logger) {
	v.Tree().log(l, v.Name())
}

func (t *DirTree) log(logger *log.Logger, path Path) {
	logger.Enter(path, func(l *log.Logger) {
		l.Println(t.Tree.String())
		for _, file := range t.Files() {
			l.Println(file.String())
		}
		for _, dir := range t.Dirs() {
			dir.log(l)
		}
	})
}

func (t *DirTree) String() string {
	var buf bytes.Buffer
	l := log.NewBufLogger(&buf)
	t.log(l, GetPath(""))
	return buf.String()
}

func (t *DirTree) Files() (files []File) {
	for it := t.Tree.Iterator(); it.Next(); {
		switch v := it.Value().(type) {
		case *Value:
			files = append(files, v.File())
		}
	}
	return
}

func (t *DirTree) Dirs() (dirs []*DirValue) {
	for it := t.Tree.Iterator(); it.Next(); {
		switch v := it.Value().(type) {
		case *DirValue:
			dirs = append(dirs, v)
		}
	}
	return
}

func PathComparator(a, b interface{}) int {
	return utils.StringComparator(a.(Path).String(), b.(Path).String())
}

func NewDirTree() *DirTree {
	return &DirTree{avltree.NewWith(PathComparator)}
}

type DirTreeValue interface {
	Name() Path
	Value() interface{}
	IsDir() bool
}

type dirTreeValue struct {
	name  Path
	value interface{}
}

func (v *dirTreeValue) Name() Path {
	return v.name.Init()
}

func (v *dirTreeValue) Value() interface{} {
	return v.value
}

func (v *dirTreeValue) IsDir() bool {
	_, ok := v.value.(*DirTree)
	return ok
}

type Value struct {
	dirTreeValue
}

func (v *Value) File() File {
	return v.value.(File)
}

func (v Value) String() string {
	return v.File().String()
}

func NewValue(file File) *Value {
	v := new(Value)
	v.name = GetPath(file.Name)
	v.value = file
	return v
}

type DirValue struct {
	dirTreeValue
}

func NewDirValue(path Path) *DirValue {
	v := new(DirValue)
	v.name = path
	v.value = NewDirTree()
	return v
}

func (v *DirValue) Tree() *DirTree {
	return v.value.(*DirTree)
}

func (v DirValue) String() string {
	return fmt.Sprintf("dir: %s", v.name)
}

type DirTreeNode avltree.Node

func dirTreeNode(node *avltree.Node) *DirTreeNode {
	return (*DirTreeNode)(node)
}

func (n *DirTreeNode) GetValue() DirTreeValue {
	return n.Value.(DirTreeValue)
}

func (n *DirTreeNode) DirValue() (*DirValue, bool) {
	dir, ok := n.Value.(*DirValue)
	return dir, ok
}

func (n *DirTreeNode) FileValue() (*Value, bool) {
	f, ok := n.Value.(*Value)
	return f, ok
}

func (t *DirTree) Get(path Path) (DirTreeValue, bool) {
	if v, ok := t.Tree.Get(path); ok {
		if val, ok := v.(DirTreeValue); ok {
			return val, true
		}
	}
	return nil, false
}

func (t *DirTree) Add(v DirTreeValue) {
	//log.Printf("adding %s", v)
	t.Put(v.Name(), v)
}

func (t *DirTree) PutFilePath(fp string, typ Type) (Path, *DirTree, *Value) {
	var base *DirTree
	path, file := GetFile(fp, typ)
	if path.IsRoot() {
		base = t
	} else if tree := t.PutDir(path); tree != nil {
		base = tree
	}
	if base != nil {
		return path, base, base.AddFile(file)
	}
	return path, nil, nil
}

func (t *DirTree) AddFile(file File) *Value {
	v := NewValue(file)
	t.Add(v)
	return v
}

func (t *DirTree) PutDir(path Path) *DirTree {
	if path.IsRoot() {
		return t
	} else if cur, ok := t.Get(path); ok {
		dir := cur.(*DirValue)
		//log.Printf("cur[%s] = %s", path, dir.String())
		return dir.Tree()
	} else if tree, ok := t.putCeil(path); ok {
		//log.Printf("ceil[%s] = %s", path, tree.String())
		return tree
	} else if tree, ok := t.putFloor(path); ok {
		//log.Printf("floor[%s] = %s", path, tree.String())
		return tree
	}
	return t.addDir(path)
}

func (t *DirTree) addDir(path Path) *DirTree {
	v := NewDirValue(path)
	t.Add(v)
	return v.Tree()
}

func (t *DirTree) putFloor(path Path) (*DirTree, bool) {
	if floor, ok := t.Floor(path); ok {
		return t.putNode(path, dirTreeNode(floor))
	}
	return nil, false
}

func (t *DirTree) putCeil(path Path) (*DirTree, bool) {
	if ceil, ok := t.Ceiling(path); ok {
		return t.putNode(path, dirTreeNode(ceil))
	}
	return nil, false
}

func (t *DirTree) putNode(path Path, node *DirTreeNode) (*DirTree, bool) {
	tree, i := node.append(t, path)
	if i == -1 {
		return tree, tree != nil
	}
	dir := tree.PutDir(path[i:])
	return dir, dir != nil
}

func (n *DirTreeNode) append(parent *DirTree, newPath Path) (*DirTree, int) {
	dir, ok := n.DirValue()
	if !ok {
		return nil, -1
	}
	path := dir.name
	i, ok := path.IndexMatch(newPath)
	if !ok {
		return nil, -1
	}

	i++
	if i == len(path) {
		return dir.Tree().PutDir(newPath[i:]), -1
	}

	// Remove node path
	parent.Remove(path)

	// Add new path as base
	oldDir := dir
	dir = NewDirValue(path[:i])
	parent.Add(dir)
	tree := dir.Tree()

	// update old path, add to base
	oldDir.name = path[i:]
	tree.Add(oldDir)
	return tree, i
}

type DirList list.List

func (dl *DirList) Front() *DirElement {
	return dirElement(dl.list().Front())
}

func (dl *DirList) Back() *DirElement {
	return dirElement(dl.list().Back())
}

func (dl *DirList) ForEach(each func(d Dir)) {
	for de := dl.Front(); de != nil; de = de.Next() {
		each(de.Get())
	}
}

func (dl *DirList) Add(dir Dir) *DirElement {
	list := dl.list()
	if list.Len() == 0 {
		e := list.PushFront(dir)
		return dirElement(e)
	}
	path := dir.Path().String()
	for de := dl.Front(); de != nil; de = de.Next() {
		if de.Get().Path().String() > path {
			e := dl.list().InsertBefore(dir, de.element())
			return dirElement(e)
		}
	}
	e := dl.list().PushBack(dir)
	return dirElement(e)
}

func (dl *DirList) Remove(dir Dir) *DirElement {
	for de := dl.Front(); de != nil; de = de.Next() {
		if de.Get() == dir {
			dl.list().Remove(de.element())
			return de
		}
	}
	return nil
}

func (dl *DirList) list() *list.List {
	return (*list.List)(dl)
}

type DirElement list.Element

func dirElement(e *list.Element) *DirElement {
	return (*DirElement)(e)
}

func (de *DirElement) Prev() *DirElement {
	return (*DirElement)(de.element().Prev())
}

func (de *DirElement) Next() *DirElement {
	return (*DirElement)(de.element().Next())
}

func (de *DirElement) Get() Dir {
	return de.Value.(Dir)
}

func (de *DirElement) element() *list.Element {
	return (*list.Element)(de)
}

type dirList []Dir
type dir struct {
	path  Path
	index int
	dirs  dirList
	files []File
}

func NewDir() Dir {
	return new(dir)
}

func (d *dir) Path() Path {
	return d.path.Init()
}

func (d *dir) Parent() Path {
	return d.Path()[:d.index]
}

func (d *dir) Base() Path {
	return d.Path()[d.index:]
}

func (d *dir) Dirs() []Dir {
	return d.dirs
}

func (d *dir) Files() []File {
	return d.files
}

func (d *dir) AddFile(f string, t Type) Dir {
	return d.Add(GetFile(f, t))
}

func (d *dir) Add(path Path, file File) Dir {
	if dir := d.addPath(path); dir != nil {
		dir.appendFile(file)
		return dir
	}
	return nil
}

func (d *dir) addPath(path Path) Dir {
	base := d.Base()
	switch {
	case base.String() == path.String():
		return d
	case base.IsRoot():
		return d.dirs.add(path, 0)
	}

	i, ok := base.IndexMatch(path)
	if !ok {
		return nil
	}

	i++
	if i == len(base) {
		return d.dirs.add(path, i)
	}

	//newD := dir{
	//	path:  d.path[:d.index+i],
	//	index: d.index,
	//}
	//parent := *d
	//oldDir.name = name[i:]
	//oldDir.parent = d
	//*d = dir{
	//	parent: d.parent,
	//	name:   name[:i],
	//}
	//
	//d.dirs.append(&oldDir)
	//return d.addPath(path)
	return nil
}

func (d *dir) appendFile(file File) {
	d.files = append(d.files, file)
}

func (d dir) MarshalJSON() ([]byte, error) {
	m := make(map[string]interface{}, len(d.dirs))
	for _, dir := range d.dirs {
		m[dir.Base().String()] = dir
	}
	if len(d.files) > 0 {
		m["."] = d.files
	}
	return json.Marshal(m)
}

func (list *dirList) add(path Path, index int) Dir {
	for _, c := range *list {
		if dir := c.addPath(path[index:]); dir != nil {
			return dir
		}
	}
	return list.append(&dir{
		path:  path,
		index: index,
	})
}

func (list *dirList) append(dir Dir) Dir {
	*list = append(*list, dir)
	return dir
}
