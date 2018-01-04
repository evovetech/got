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
	"github.com/evovetech/got/log"
)

type Module interface {
	Dir

	Src() Dir
}

type module struct {
	*dir
}

func createModule(d Dir) (Module, bool) {
	if m, ok := d.(*dir); ok {
		return &module{m}, true
	}
	return nil, false
}

func (m *module) Src() Dir {
	path := SrcPath()
	if src, ok := m.GetDir(path); ok {
		return src
	}
	return m.PutDir(path)
}

func (m *module) Copy() Entry {
	if d, ok := m.dir.Copy().(*dir); ok {
		return &module{d}
	}
	return nil
}

func (m *module) String() string {
	return DirString(m)
}

func (m *module) log(l *log.Logger) {
	logDir(l, "module", m)
}

//
//func (m *module) allModules() (modules []*module) {
//	for _, temp := range m.Modules() {
//		mod := temp.(*module)
//		modules = append(modules, mod)
//		for _, child := range mod.allModules() {
//			cp := *child
//			cp.setPath(cp.Path().CopyWithParent(mod.Path()))
//
//		}
//	}
//}
