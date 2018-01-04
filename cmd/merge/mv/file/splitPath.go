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

import "fmt"

type SplitPath interface {
	Val() Path
	Len() int
	Next() SplitPath
}

func NewSplitPath(path Path, length int) SplitPath {
	return splitPath{path, length}
}

func (p Path) splitAt(index int) SplitPath {
	return NewSplitPath(p, index)
}

func (p Path) split() SplitPath {
	return p.splitAt(len(p))
}

type splitPath struct {
	path Path
	len  int
}

func (p splitPath) Val() Path {
	return p.path[:p.len]
}

func (p splitPath) Len() int {
	return p.len
}

func (p splitPath) Next() SplitPath {
	return p.path[p.len:].split()
}

func (p splitPath) String() string {
	return fmt.Sprintf("'%s'|'%s'", p.Val(), p.Next().Val())
}
