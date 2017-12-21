package file

import (
	"fmt"
	"github.com/evovetech/got/log"
)

type Part interface {
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
