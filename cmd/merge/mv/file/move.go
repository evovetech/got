package file

import (
	"fmt"
	"github.com/evovetech/got/log"
	"github.com/evovetech/got/util"
	"regexp"
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

type MovePath []MovePart

var reMovePart = regexp.MustCompile("\\((.*)\\)->\\((.*)\\)")

func ParseMovePath(path Path) (MovePath, bool) {
	var mv MovePath
	for _, p := range path {
		var m MovePart
		if match := reMovePart.FindStringSubmatch(p); match != nil {
			from := strings.Replace(match[1], "|", "/", -1)
			to := strings.Replace(match[2], "|", "/", -1)
			m = NewUnequalMovePart(from, to)
		} else {
			m = NewEqualMovePart(p)
		}
		if !mv.Append(m) {
			return nil, false
		}
	}
	return mv, true
}

func (mv *MovePath) Append(parts ...MovePart) (added bool) {
	for _, m := range parts {
		if m.Max() > 0 {
			added = true
			*mv = append(*mv, m)
		}
	}
	return
}

func (mv MovePath) FromPath() Path {
	return mv.makePath(func(part MovePart) Path {
		return part.From().Val()
	})
}

func (mv MovePath) ToPath() Path {
	return mv.makePath(func(part MovePart) Path {
		return part.To().Val()
	})
}

func (mv MovePath) Path() Path {
	return mv.makePath(func(part MovePart) Path {
		return part.Path()
	})
}

func (mv MovePath) String() string {
	return mv.Path().String()
}

func (mv MovePath) makePath(f func(MovePart) Path) Path {
	var paths = make([]Path, len(mv))
	for i, part := range mv {
		paths[i] = f(part)
	}
	return JoinPaths(paths...)
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
