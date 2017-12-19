package file

import (
	"github.com/emirpasic/gods/trees/avltree"
	"github.com/evovetech/got/log"
	"reflect"
)

var (
	entryType  = reflect.TypeOf((*Entry)(nil)).Elem()
	fileType   = reflect.TypeOf((*File)(nil)).Elem()
	dirType    = reflect.TypeOf((*Dir)(nil)).Elem()
	moduleType = reflect.TypeOf((*Module)(nil)).Elem()
)

func (d *dir) Get(path Path) (Entry, bool) {
	if path.IsRoot() {
		return d, true
	}
	if v, ok := d.tree().Get(path); ok {
		if val, ok := v.(Entry); ok {
			return val, true
		}
	}
	return nil, false
}

func (d *dir) GetDir(path Path) (Dir, bool) {
	if e, ok := d.Get(path); ok {
		return e.Dir()
	}
	return nil, false
}

func (d *dir) Find(path Path) (Entry, bool) {
	if e, ok := d.Get(path); ok {
		return e, ok
	}

	for _, finder := range []*NodeFinder{d.Floor(), d.Ceiling()} {
		if n, ok := finder.Get(path); ok {
			if dir, found := n.find(path); found {
				return dir, ok
			}
		}
	}

	return nil, false
}

func (d *dir) FindDir(path Path) (Dir, bool) {
	if e, ok := d.Find(path); ok {
		return e.Dir()
	}
	return nil, false
}

func (d *dir) PutFile(fp string, typ Type) (Dir, File) {
	path, f := GetFile(fp, typ)
	if dir := d.PutDir(path); dir != nil {
		dir.put(f)
		return dir, f
	}
	return nil, nil
}

func (d *dir) PutDir(path Path) Dir {
	if dir, ok := d.putModule(path); ok {
		return dir
	}
	return d.insertDir(path)
}

func trueFilter(t reflect.Type) bool {
	return true
}

func (d *dir) filter(t reflect.Type, f EntryFilter) interface{} {
	st := reflect.SliceOf(t)
	slice := reflect.MakeSlice(st, 0, 0)
	for it := d.tree().Iterator(); it.Next(); {
		v := reflect.ValueOf(it.Value())
		typ := v.Type()
		if typ.AssignableTo(t) && f(typ) {
			slice = reflect.Append(slice, v)
		}
	}
	return slice.Interface()
}

func (d *dir) collect(t reflect.Type) interface{} {
	return d.filter(t, trueFilter)
}

func (d *dir) Entries() []Entry {
	return d.collect(entryType).([]Entry)
}

func (d *dir) Files() []File {
	return d.collect(fileType).([]File)
}

func (d *dir) Dirs() []Dir {
	return d.filter(dirType, func(t reflect.Type) bool {
		return !t.AssignableTo(moduleType)
	}).([]Dir)
}

func (d *dir) Modules() []Module {
	return d.collect(moduleType).([]Module)
}

func (d *dir) AllEntries(filter EntryFilter) (entries []Entry) {
	for _, e := range d.filter(entryType, filter).([]Entry) {
		entries = append(entries, e)
		if dir, ok := e.(Dir); ok {
			entries = append(entries, dir.AllEntries(filter)...)
		}
	}
	return entries
}

func (d *dir) AllFiles() []File {
	files := d.Files()
	for it := d.tree().Iterator(); it.Next(); {
		reflect.TypeOf(it.Value())
		switch dir := it.Value().(type) {
		case Dir:
			for _, f := range dir.AllFiles() {
				files = append(files, f.CopyWithParent(dir.Path()))
			}
		}
	}
	return files
}

func (d *dir) MvCount() TypeCount {
	tc := make(TypeCount)
	for _, f := range d.Files() {
		tc.add(f.Type(), 1)
	}
	for _, dir := range d.Dirs() {
		tc.addAll(dir.MvCount())
	}
	return tc
}

func (d *dir) insertDir(path Path) Dir {
	if path.IsRoot() {
		return d
	}

	for _, finder := range []*NodeFinder{d.Floor(), d.Ceiling()} {
		if dir, ok := finder.put(path); ok {
			return dir
		}
	}

	return d.putDir(path)
}

func (d *dir) putModule(path Path) (Dir, bool) {
	if path.IsRoot() {
		return d, true
	} else if i, ok := path.SrcIndex(); ok {
		if mod := d.getModule(path[:i]); mod != nil {
			return mod.insertDir(path[i:]), true
		}
	}
	return nil, false
}

func (d *dir) getModule(path Path) (mod Module) {
	var ok bool
	var parent Dir
	var dir = d.PutDir(path)
	if mod, ok = dir.(Module); !ok {
		log.Std.Enter(d.Path(), func(l *log.Logger) {
			l.Printf("path=%s, mod=%s", path, dir.Path())
			if i := len(path) - len(dir.Path()); i > 0 {
				parent, ok = d.FindDir(path[:i])
			} else {
				parent = d
			}
			if parent != nil {
				l.Printf("mod parent: %s", parent.Path())
				if mod, ok = createModule(dir); ok {
					parent.put(mod)
				}
			}
		})
	}
	return
}

func (d *dir) put(e Entry) {
	d.tree().Put(e.Path(), e)
}

func (d *dir) putDir(path Path) Dir {
	e := NewDir(path)
	d.put(e)
	return e
}

type NodeFinderFunc func(key interface{}) (node *avltree.Node, found bool)
type NodeFinder struct {
	d    *dir
	call NodeFinderFunc
}

func (find *NodeFinder) Get(path Path) (n *Node, found bool) {
	if tn, f := find.call(path); f {
		n, found = node(tn), f
	}
	return
}

func (find *NodeFinder) put(path Path) (Dir, bool) {
	if n, ok := find.Get(path); ok {
		return find.d.putNode(path, n)
	}
	return nil, false
}

func (d *dir) Floor() *NodeFinder {
	return &NodeFinder{d, d.tree().Floor}
}

func (d *dir) Ceiling() *NodeFinder {
	return &NodeFinder{d, d.tree().Ceiling}
}

func (d *dir) putNode(path Path, node *Node) (Dir, bool) {
	if e := node.Entry(); path.Equals(e.Path()) {
		return e.(Dir), true
	}
	if tree, i := node.append(d, path); i > 0 {
		dir := tree.PutDir(path[i:])
		return dir, dir != nil
	} else {
		return tree, tree != nil
	}
}
