package tree

import (
	"fmt"
	"github.com/evovetech/got/git/object"
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

type Tree interface {
	object.Object

	List() EntryList
}

type tree struct {
	object.Object

	list EntryList
}

func New(id object.Id) Tree {
	t := &tree{
		Object: object.NewTree(id),
	}
	t.SetInitFunc(func() {
		t.list = Ls(t.Id())
	})
	return t
}

func (t *tree) List() (l EntryList) {
	return t.init().list
}

func (t *tree) String() string {
	return fmt.Sprintf(treeFormat, t.Id(), t.List())
}

func (t *tree) init() *tree {
	t.Object.Init()
	return t
}
