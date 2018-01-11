package object

import (
	"fmt"
	"github.com/evovetech/got/git/types"
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

	list ObjectList
}

func NewTree(id types.Id) Tree {
	t := &tree{
		Object: New(id, types.Tree),
	}
	t.SetInitFunc(func() {
		t.list = Ls(t.Id())
	})
	return t
}

func (t *tree) List() (l ObjectList) {
	return t.init().list
}

func (t *tree) String() string {
	return fmt.Sprintf(treeFormat, t.Id(), t.List())
}

func (t *tree) init() *tree {
	t.Object.Init()
	return t
}
