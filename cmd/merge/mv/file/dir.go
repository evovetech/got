package file

import (
	"container/list"
	"encoding/json"
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
	//oldDir.path = path[i:]
	//oldDir.parent = d
	//*d = dir{
	//	parent: d.parent,
	//	path:   path[:i],
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
