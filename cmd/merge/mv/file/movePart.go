package file

import (
	"fmt"
)

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
