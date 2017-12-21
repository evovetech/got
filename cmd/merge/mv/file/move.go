package file

import (
	"fmt"
	"github.com/evovetech/got/log"
	"github.com/evovetech/got/util"
	"strings"
)

type Move interface {
	fmt.Stringer
	Key() Path
	Value() interface{}
	Path() Path
	IsDir() bool

	// private
	log(l *log.Logger)
}

type move struct {
}

type MovePath []MovePart

func (mv *MovePath) Append(parts ...MovePart) (added bool) {
	for _, m := range parts {
		if m.Max() > 0 {
			added = true
			*mv = append(*mv, m)
		}
	}
	return
}

func (mv MovePath) Path() Path {
	var paths []Path
	for _, part := range mv {
		paths = append(paths, part.Path())
	}
	return JoinPaths(paths...)
}

func (mv MovePath) String() string {
	return mv.Path().String()
}

func (mv MovePath) log(l *log.Logger) {
	var result [5][]string
	for _, m := range mv {
		var p [5]string
		var max int
		if m.Equal() {
			s := m.From().Val().String()
			max = len(s)
			p[2] = s
		} else {
			f := m.From().Val().String()
			t := m.To().Val().String()
			max = util.MaxInt(len(f), len(t))
			p[1], p[3] = f, t
		}
		for i, str := range p {
			sym := '-'
			if i == 0 || i == 4 {
				sym = '*'
			}
			padEnd(&str, max, sym)
			result[i] = append(result[i], str)
		}
	}
	for _, r := range result {
		l.Println(strings.Join(r, "|"))
	}
}

func (mv MovePath) get() (MovePath, bool) {
	return mv, len(mv) > 0
}

func padEnd(s *string, max int, r rune) {
	pad := string(r)
	if size := len(*s); max > size {
		p := max - size
		for i := 0; i < p; i++ {
			*s = *s + pad
		}
	}
}
