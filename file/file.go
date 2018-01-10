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
	"encoding/json"
	"fmt"
	"github.com/evovetech/got/log"
)

type File interface {
	Entry
	Name() string
	Type() Type
	CopyWithParent(parent Path) File
}

type file struct {
	entry
}

func GetFile(file string, typ Type) (Path, File) {
	path := GetPath(file)
	return path.Dir(), NewFile(path.Name(), typ)
}

func ReverseFile(file string, typ Type) (Path, File) {
	path := GetPath(file)
	path.Reverse()
	return path.Dir(), NewFile(path.Name(), typ)
}

func NewFile(file string, typ Type) File {
	return NewFileWithPath(GetPath(file), typ)
}

func NewFileWithPath(path Path, typ Type) File {
	f := new(file)
	f.path = path
	f.value = typ
	return f
}

func (f *file) Name() string {
	return f.Key().Name()
}

func (f *file) Type() Type {
	return f.value.(Type)
}

func (f *file) IsDir() bool {
	return false
}

func (f *file) Copy() Entry {
	return NewFileWithPath(f.path.Copy(), f.Type())
}

func (f *file) CopyWithParent(parent Path) File {
	path := f.Key().CopyWithParent(parent)
	return NewFileWithPath(path, f.Type())
}

func (f file) String() string {
	return fmt.Sprintf("%s: '%s'", f.Type(), f.Key().String())
}

func (f file) MarshalJSON() ([]byte, error) {
	return json.Marshal(f.String())
}

func (f *file) log(l *log.Logger) {
	l.Println(f.String())
}
