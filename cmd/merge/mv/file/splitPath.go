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
