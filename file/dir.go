/*
 * Copyright 2018 evove.tech
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package file

import (
	"bytes"
	"fmt"
	"github.com/emirpasic/gods/trees/avltree"
	"github.com/evovetech/got/log"
	"reflect"
)

type EntryFilter func(reflect.Type) bool

type Dir interface {
	Entry

	Get(path Path) (Entry, bool)
	GetDir(path Path) (Dir, bool)
	Find(path Path) (Entry, bool)
	FindDir(path Path) (Dir, bool)
	PutFile(fp string, typ Type) (Dir, File)
	PutDir(path Path) Dir

	Entries() []Entry
	Files() []File
	Dirs() []Dir
	Modules() []Module

	DeepIterator() DeepIterator
	AllEntries() []Entry
	AllFiles() []File
	AllDirs() []Dir
	AllModules() []Module

	MvCount() TypeCount

	// private
	tree() *avltree.Tree
	insertDir(path Path) Dir
	put(e Entry)
	putDir(path Path) Dir
}

type dir struct {
	entry
}

func NewRoot() Dir {
	return NewDir(GetPath(""))
}

func NewDir(path Path) Dir {
	e := new(dir)
	e.path = path
	e.value = avltree.NewWith(PathComparator)
	return e
}

func DirString(d Dir) string {
	var buf bytes.Buffer
	l := log.NewBufLogger(&buf)
	d.log(l)
	return buf.String()
}

func (d *dir) Copy() Entry {
	if d == nil {
		return nil
	}
	cp := NewDir(d.Key())
	for _, e := range d.Entries() {
		cp.put(e.Copy())
	}
	return cp
}

func (d *dir) IsDir() bool {
	return true
}

func (d *dir) Iterator() Iterator {
	return newIterator(d.tree())
}

func (d *dir) DeepIterator() DeepIterator {
	return newDeepIterator(d)
}

func (d *dir) String() string {
	return DirString(d)
}

func (d *dir) tree() *avltree.Tree {
	return d.value.(*avltree.Tree)
}

func logDir(l *log.Logger, name string, d Dir) {
	prefix := fmt.Sprintf("%s<%s>", name, d.Key().String())
	l.Enter(prefix, func(_ *log.Logger) {
		for t, n := range d.MvCount() {
			l.Printf("%s: %d\n", t.String(), n)
		}
		for _, e := range d.Entries() {
			e.log(l)
		}
	})
}

func logDir2(l *log.Logger, name string, d Dir) {
	prefix := fmt.Sprintf("%s<%s>", name, d.Key().String())
	if files := d.Files(); len(files) > 0 {
		l.Enter(prefix, func(_ *log.Logger) {
			for t, n := range d.MvCount() {
				l.Printf("%s: %d\n", t.String(), n)
			}
			for _, f := range files {
				f.log(l)
			}
		})
	} else {
		l.Println(prefix)
	}
	for _, e := range d.AllDirs() {
		e.log(l)
	}
}

func (d *dir) log(l *log.Logger) {
	logDir(l, "dir", d)
}
