package tree

import (
	"fmt"
	"github.com/evovetech/got/git"
	"github.com/evovetech/got/git/object"
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
	git.Object

	list git.ObjectList
}

func New(id types.Id) git.Tree {
	t := &tree{
		Object: object.NewTree(id),
	}
	t.SetInitFunc(func() {
		t.list = Ls(t.Id())
	})
	return t
}

func (t *tree) List() (l git.ObjectList) {
	return t.init().list
}

func (t *tree) String() string {
	return fmt.Sprintf(treeFormat, t.Id(), t.List())
}

func (t *tree) init() *tree {
	t.Object.Init()
	return t
}
