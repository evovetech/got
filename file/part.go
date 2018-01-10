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
	"github.com/evovetech/got/util"
)

type Part interface {
	fmt.Stringer
	From() SplitPath
	To() SplitPath
	Max() int

	next() Part
	matchEqual() (MovePart, bool)
	matchUnequal() (MovePart, bool)
}

type partMatch func() (MovePart, bool)

type part struct {
	from SplitPath
	to   SplitPath
}

func NewPart(from SplitPath, to SplitPath) Part {
	return part{from, to}
}

func (p part) From() SplitPath {
	return p.from
}

func (p part) To() SplitPath {
	return p.to
}

func (p part) Max() int {
	return util.MaxInt(p.from.Len(), p.to.Len())
}

func (p part) next() Part {
	from := p.from.Next()
	to := p.to.Next()
	return part{from, to}
}

func (p part) String() string {
	return fmt.Sprintf("part {\n  from: %s,\n  to: %s\n}", p.from, p.to)
}

func (p part) matchEqual() (MovePart, bool) {
	l := log.Verbose
	l.Println("matchEqual")
	l.In()
	defer l.Out()
	l.Printf("try: %s", p)
	from := p.From().Val()
	to := p.To().Val()
	if i, ok := from.IndexMatch(to); ok {
		m := NewMovePart(
			from.splitAt(i+1),
			to.splitAt(i+1),
			true,
		)
		l.Printf("found: %s", m)
		return m, true
	}
	return nil, false
}

func (p part) matchUnequal() (MovePart, bool) {
	l := log.Verbose
	l.Println("matchUnequal")
	l.In()
	defer l.Out()
	l.Printf("try: %s", p)
	from := p.From().Val()
	to := p.To().Val()
	if f, t, ok := nextMatch(from, to); ok {
		m := NewMovePart(
			from.splitAt(f),
			to.splitAt(t),
			false,
		)
		l.Printf("found: %s", m)
		return m, true
	}
	return nil, false
}
