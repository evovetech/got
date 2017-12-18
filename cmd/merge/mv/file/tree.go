package file

func (d *dir) Get(path Path) (Entry, bool) {
	if v, ok := d.tree().Get(path); ok {
		if val, ok := v.(Entry); ok {
			return val, true
		}
	}
	return nil, false
}

func (d *dir) PutFile(fp string, typ Type) (Dir, File) {
	path, f := GetFile(fp, typ)
	if dir := d.PutDir(path); dir != nil {
		dir.add(f)
		return dir, f
	}
	return nil, nil
}

func (d *dir) PutDir(path Path) Dir {
	if path.IsRoot() {
		return d
	} else if tree, ok := d.putFloor(path); ok {
		return tree
	} else if tree, ok := d.putCeil(path); ok {
		return tree
	}
	return d.addDir(path)
}

func (d *dir) Files() (files []File) {
	for it := d.tree().Iterator(); it.Next(); {
		switch e := it.Value().(type) {
		case File:
			files = append(files, e)
		}
	}
	return
}

func (d *dir) Dirs() (dirs []Dir) {
	for it := d.tree().Iterator(); it.Next(); {
		switch e := it.Value().(type) {
		case Dir:
			dirs = append(dirs, e)
		}
	}
	return
}

type TypeCount map[Type]int

func (c TypeCount) add(t Type, i int) {
	c[t] = c[t] + i
}

func (c TypeCount) addAll(o TypeCount) {
	for t, n := range o {
		c.add(t, n)
	}
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

func (d *dir) add(e Entry) {
	d.tree().Put(e.Path(), e)
}

func (d *dir) addDir(path Path) Dir {
	e := NewDir(path)
	d.add(e)
	return e
}

func (d *dir) putFloor(path Path) (Dir, bool) {
	if floor, ok := d.tree().Floor(path); ok {
		return d.putNode(path, node(floor))
	}
	return nil, false
}

func (d *dir) putCeil(path Path) (Dir, bool) {
	if ceil, ok := d.tree().Ceiling(path); ok {
		return d.putNode(path, node(ceil))
	}
	return nil, false
}

func (d *dir) putNode(path Path, node *Node) (Dir, bool) {
	if e := node.Entry(); path.Equals(e.Path()) {
		return e.(Dir), true
	}
	tree, i := node.append(d, path)
	if i == -1 {
		return tree, tree != nil
	}
	dir := tree.PutDir(path[i:])
	return dir, dir != nil
}
