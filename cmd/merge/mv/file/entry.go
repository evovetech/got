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
	"fmt"
	"github.com/evovetech/got/log"
)

type Entry interface {
	fmt.Stringer
	Key() Path
	Value() interface{}
	Path() Path
	IsDir() bool
	File() (File, bool)
	Dir() (Dir, bool)
	Copy() Entry
	Iterator() Iterator

	// private
	setPath(path Path)
	log(l *log.Logger)
}

type entry struct {
	path  Path
	value interface{}
}

func (e *entry) Value() interface{} {
	return e.value
}

func (e *entry) Key() Path {
	return e.Path()
}

func (e *entry) Path() Path {
	return e.path.Init()
}

func (e *entry) File() (File, bool) {
	f, ok := e.value.(File)
	return f, ok
}

func (e *entry) Dir() (Dir, bool) {
	d, ok := e.value.(Dir)
	return d, ok
}

func (e *entry) Iterator() Iterator {
	return noEntries
}

func (e *entry) setPath(path Path) {
	e.path = path
}
