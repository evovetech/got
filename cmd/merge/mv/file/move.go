package file

import (
	"bytes"
	"fmt"
	"github.com/evovetech/got/log"
	"github.com/evovetech/got/util"
	"strings"
)

type Move interface {
	Entry
	From() Rename
	To() Rename
}

type move struct {
	*dir
	from rename
	to   rename
}

type Rename interface {
	File() File
	Dir() Path
}

type rename struct {
	File
	path Path
}

func newRename(file string, typ Type) *rename {
	r := new(rename)
	r.path, r.File = GetFile(file, typ)
	return r
}

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
	return p.path.Len()
}

func (p splitPath) Next() SplitPath {
	return p.path[p.len:].split()
}

func (p splitPath) String() string {
	return fmt.Sprintf("'%s' - '%s'", p.Val(), p.Next().Val())
}

type MoveParts []MovePart

func (mp *MoveParts) Append(parts ...MovePart) (added bool) {
	for _, m := range parts {
		if m.Max() > 0 {
			added = true
			*mp = append(*mp, m)
		}
	}
	return
}

func padEnd(s *string, max int) {
	if size := len(*s); max > size {
		pad := max - size
		for i := 0; i < pad; i++ {
			*s = *s + " "
		}
	}
}

func (mp MoveParts) String() string {
	var result [3][]string
	for _, m := range mp {
		var p [3]string
		var max int
		if m.Equal() {
			s := m.From().Val().String()
			max = len(s)
			p[1] = s
		} else {
			f := m.From().Val().String()
			t := m.To().Val().String()
			max = util.MaxInt(len(f), len(t))
			p[0], p[2] = f, t
		}
		for i, str := range p {
			padEnd(&str, max)
			result[i] = append(result[i], str)
		}
	}
	var buf bytes.Buffer
	l := log.NewBufLogger(&buf)
	l.Enter("move", func(_ *log.Logger) {
		for _, r := range result {
			l.Println(strings.Join(r, "/"))
		}
	})
	return buf.String()
}

func (mp MoveParts) get() (MoveParts, bool) {
	return mp, len(mp) > 0
}

type Part interface {
	From() SplitPath
	To() SplitPath
	Max() int

	next() Part
	matchEqual() (MovePart, bool)
	matchUnequal() (MovePart, bool)
}

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
	max := p.from.Len()
	if t := p.to.Len(); t > max {
		max = t
	}
	return max
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

func (p Path) nextMatch(o Path) (int, int, bool) {
	max := util.MinInt(len(p), len(o))
	for pi := 0; pi < max; pi++ {
		oMax := util.MinInt(pi+1, max)
		oi := oMax - 1
		for i := 0; i < pi; i++ {
			//log.Printf("trying i=%d, oi=%d", i, oi)
			if p[i] == o[oi] {
				return i, oi, true
			}
		}
		for oi := 0; oi < oMax; oi++ {
			//log.Printf("trying pi=%d, oi=%d", pi, oi)
			if p[pi] == o[oi] {
				return pi, oi, true
			}
		}
	}
	return -1, -1, false
}

func nextMatch(from Path, to Path) (f int, t int, ok bool) {
	if len(from) > len(to) {
		f, t, ok = from.nextMatch(to)
	} else {
		t, f, ok = to.nextMatch(from)
	}
	return
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

type MovePart interface {
	Part
	Equal() bool
}

type movePart struct {
	part
	equal bool
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

func (m movePart) String() string {
	return fmt.Sprintf("movePart {\n  from: %s,\n  to: %s,\n  equal: %s\n}", m.from, m.to, m.equal)
}

type MoveParser struct {
	from   Path
	to     Path
	result MoveParts
}

func NewMoveParser(from Path, to Path) *MoveParser {
	return &MoveParser{
		from: from,
		to:   to,
	}
}

func (p *MoveParser) Parse() (MoveParts, bool) {
	for it := p.iterator(); it.hasNext(); {
		p.result.Append(it.get()...)
	}
	return p.result.get()
}

func (p *MoveParser) iterator() *moveIterator {
	return &moveIterator{
		cur: NewPart(p.from.split(), p.to.split()),
	}
}

type moveIterator struct {
	cur  Part
	done bool
}

func (it *moveIterator) hasNext() bool {
	return !it.done
}

func (it *moveIterator) get() (parts MoveParts) {
	done := func() {
		if last := it.cur; last.Max() > 0 {
			parts.Append(NewMovePart(
				last.From(),
				last.To(),
				false),
			)
		}
		it.cur = nil
		it.done = true
	}

	if m, ok := it.cur.matchEqual(); ok && parts.Append(m) {
		it.cur = m.next()
	} else {
		done()
		return
	}

	if m, ok := it.cur.matchUnequal(); ok && parts.Append(m) {
		it.cur = m.next()
	} else {
		done()
	}
	return
}

//
//func (d *dir) AddMove(from string, to string) Move {
//	if !d.Path().IsRoot() {
//		return nil
//	}
//
//	//fromPath := GetPath(from)
//	//toPath := GetPath(to)
//
//	//path, f1 := GetFile(from, Rn|Del)
//	//root := NewDir(path)
//	//root.put(f1)
//	//dTo, fTo := root.PutFile(to, Rn|Add)
//
//	//rf := newRename(from, Rn|Del)
//	//rt := newRename(to, Rn|Add)
//	return nil
//}
