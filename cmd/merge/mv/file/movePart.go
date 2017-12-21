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
	path := fmt.Sprintf("{|%s|->|%s|}", from, to)
	return []string{path}
}

func (m movePart) String() string {
	return m.Path().String()

}

func (m movePart) log(l *log.Logger) {
	l.Printf("movePart {\n  from: %s,\n  to: %s,\n  equal: %s\n}", m.from, m.to, m.equal)
}
