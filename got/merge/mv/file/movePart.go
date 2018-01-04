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
	"strings"
)

type MovePart interface {
	Part
	Equal() bool
	Path() Path

	// private
	log(l *log.Logger)
}

type movePart struct {
	part
	equal bool
}

func NewUnequalMovePart(from string, to string) MovePart {
	return NewMovePart(
		GetPath(from).split(),
		GetPath(to).split(),
		false,
	)
}

func NewEqualMovePart(path string) MovePart {
	sp := GetPath(path).split()
	return NewMovePart(sp, sp, true)
}

func NewMovePart(from SplitPath, to SplitPath, equal bool) MovePart {
	p := new(movePart)
	p.from = from
	p.to = to
	p.equal = equal
	return p
}

func (m movePart) Equal() bool {
	return m.equal
}

func (m movePart) Path() Path {
	f := m.from.Val()
	if m.equal {
		return f
	}

	// unique path segment
	from := strings.Join(f, "|")
	to := strings.Join(m.to.Val(), "|")
	path := fmt.Sprintf("(%s)->(%s)", from, to)
	return []string{path}
}

func (m movePart) String() string {
	return m.Path().String()

}

func (m movePart) log(l *log.Logger) {
	l.Printf("movePart {\n  from: %s,\n  to: %s,\n  equal: %s\n}", m.from, m.to, m.equal)
}
