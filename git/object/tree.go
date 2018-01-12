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
	Object

	list List
}

func NewTree(id Id) Tree {
	t := &tree{
		Object: New(id, TreeType),
	}
	t.setInitFunc(func() {
		t.list = lsTree(t.Id())
	})
	return t
}

func (t *tree) List() (l List) {
	return t.init().list
}

func (t *tree) String() string {
	return fmt.Sprintf(treeFormat, t.Id(), t.List())
}

func (t *tree) init() *tree {
	t.Object.Init()
	return t
}
