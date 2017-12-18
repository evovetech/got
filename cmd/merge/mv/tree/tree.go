package tree

import (
	"github.com/evovetech/got/cmd/merge/mv/file"
)

func (de *dirEntry) Path() file.Path {
	return de.Key()
}

func (de *dirEntry) Get(path file.Path) (Entry, bool) {
	if v, ok := de.tree().Get(path); ok {
		if val, ok := v.(Entry); ok {
			return val, true
		}
	}
	return nil, false
}

func (de *dirEntry) PutFile(fp string, typ file.Type) (DirEntry, FileEntry) {
	path, f := file.GetFile(fp, typ)
	if dir := de.PutDir(path); dir != nil {
		return dir, dir.addFile(f)
	}
	return nil, nil
}

func (de *dirEntry) PutDir(path file.Path) DirEntry {
	if path.IsRoot() {
		return de
	} else if tree, ok := de.putFloor(path); ok {
		return tree
	} else if tree, ok := de.putCeil(path); ok {
		return tree
	}
	return de.addDir(path)
}

func (de *dirEntry) add(e Entry) {
	de.tree().Put(e.Key(), e)
}

func (de *dirEntry) addFile(file file.File) FileEntry {
	e := NewFileEntry(file)
	de.add(e)
	return e
}

func (de *dirEntry) addDir(path file.Path) DirEntry {
	e := NewDirEntry(path)
	de.add(e)
	return e
}

func (de *dirEntry) Files() (files []FileEntry) {
	for it := de.tree().Iterator(); it.Next(); {
		switch e := it.Value().(type) {
		case FileEntry:
			files = append(files, e)
		}
	}
	return
}

func (de *dirEntry) Dirs() (dirs []DirEntry) {
	for it := de.tree().Iterator(); it.Next(); {
		switch e := it.Value().(type) {
		case DirEntry:
			dirs = append(dirs, e)
		}
	}
	return
}

func (de *dirEntry) MvCount() (add int, del int) {
	for _, e := range de.Files() {
		f := e.File()
		switch {
		case f.Type.HasFlag(file.Add):
			add++
		case f.Type.HasFlag(file.Del):
			del++
		}
	}
	for _, dir := range de.Dirs() {
		a, d := dir.MvCount()
		add += a
		del += d
	}
	return
}

func (de *dirEntry) putFloor(path file.Path) (DirEntry, bool) {
	if floor, ok := de.tree().Floor(path); ok {
		return de.putNode(path, node(floor))
	}
	return nil, false
}

func (de *dirEntry) putCeil(path file.Path) (DirEntry, bool) {
	if ceil, ok := de.tree().Ceiling(path); ok {
		return de.putNode(path, node(ceil))
	}
	return nil, false
}

func (de *dirEntry) putNode(path file.Path, node *Node) (DirEntry, bool) {
	if e := node.Entry(); path.Equals(e.Key()) {
		return e.(DirEntry), true
	}
	tree, i := node.append(de, path)
	if i == -1 {
		return tree, tree != nil
	}
	dir := tree.PutDir(path[i:])
	return dir, dir != nil
}
