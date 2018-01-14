package object

import (
	"fmt"
	"regexp"
)

var (
	reTreeLine = regexp.MustCompile("^(\\d+) (\\w+) (\\w+)\\t(.*)$")
)

const treeFormat = `{
  id: "%s",
  list: %s,
}
`

type tree struct {
	object

	list List
}

func NewTree(id Id) Tree {
	t := new(tree)
	t.id, t.kind = id, TreeType
	t.initFunc = func() {
		t.list = lsTree(t.Id())
	}
	return t
}

func (t *tree) List() (l List) {
	return t.init().list
}

func (t *tree) String() string {
	return fmt.Sprintf(treeFormat, t.Id(), t.List())
}

func (t *tree) init() *tree {
	t.object.Init()
	return t
}
