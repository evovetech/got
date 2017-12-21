package file

import (
	"bytes"
	"github.com/evovetech/got/log"
	"github.com/evovetech/got/util"
	"strings"
)

type Move []MovePart

func (mv *Move) Append(parts ...MovePart) (added bool) {
	for _, m := range parts {
		if m.Max() > 0 {
			added = true
			*mv = append(*mv, m)
		}
	}
	return
}

func (mv Move) String() string {
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
	var buf bytes.Buffer
	l := log.NewBufLogger(&buf)
	for _, r := range result {
		l.Println(strings.Join(r, "|"))
	}
	return buf.String()
}

func (mv Move) get() (Move, bool) {
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
