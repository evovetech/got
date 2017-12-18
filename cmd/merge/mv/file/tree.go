package file

func (d *dirEntry) Get(path Path) (Entry, bool) {
	if v, ok := d.tree().Get(path); ok {
		if val, ok := v.(Entry); ok {
			return val, true
		}
	}
	return nil, false
}

func (d *dirEntry) PutFile(fp string, typ Type) (DirEntry, FileEntry) {
	path, f := GetFile(fp, typ)
	if dir := d.PutDir(path); dir != nil {
		return dir, dir.addFile(f)
	}
	return nil, nil
}

func (d *dirEntry) PutDir(path Path) DirEntry {
	if path.IsRoot() {
		return d
	} else if tree, ok := d.putFloor(path); ok {
		return tree
	} else if tree, ok := d.putCeil(path); ok {
		return tree
	}
	return d.addDir(path)
}

func (d *dirEntry) Files() (files []FileEntry) {
	for it := d.tree().Iterator(); it.Next(); {
		switch e := it.Value().(type) {
		case FileEntry:
			files = append(files, e)
		}
	}
	return
}

func (d *dirEntry) Dirs() (dirs []DirEntry) {
	for it := d.tree().Iterator(); it.Next(); {
		switch e := it.Value().(type) {
		case DirEntry:
			dirs = append(dirs, e)
		}
	}
	return
}

func (d *dirEntry) MvCount() (add int, del int) {
	for _, e := range d.Files() {
		f := e.File()
		switch {
		case f.Type.HasFlag(Add):
			add++
		case f.Type.HasFlag(Del):
			del++
		}
	}
	for _, dir := range d.Dirs() {
		a, d := dir.MvCount()
		add += a
		del += d
	}
	return
}

func (d *dirEntry) add(e Entry) {
	d.tree().Put(e.Key(), e)
}

func (d *dirEntry) addFile(file File) FileEntry {
	e := NewFileEntry(file)
	d.add(e)
	return e
}

func (d *dirEntry) addDir(path Path) DirEntry {
	e := NewDirEntry(path)
	d.add(e)
	return e
}

func (d *dirEntry) putFloor(path Path) (DirEntry, bool) {
	if floor, ok := d.tree().Floor(path); ok {
		return d.putNode(path, node(floor))
	}
	return nil, false
}

func (d *dirEntry) putCeil(path Path) (DirEntry, bool) {
	if ceil, ok := d.tree().Ceiling(path); ok {
		return d.putNode(path, node(ceil))
	}
	return nil, false
}

func (d *dirEntry) putNode(path Path, node *Node) (DirEntry, bool) {
	if e := node.Entry(); path.Equals(e.Key()) {
		return e.(DirEntry), true
	}
	tree, i := node.append(d, path)
	if i == -1 {
		return tree, tree != nil
	}
	dir := tree.PutDir(path[i:])
	return dir, dir != nil
}
