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

type Tree struct {
	*object.Object

	list EntryList
}

func New(id object.Id) *Tree {
	t := &Tree{
		Object: object.NewTree(id),
	}
	t.SetInitFunc(t.init)
	return t
}

func (t *Tree) Init() *Tree {
	t.Object.Init()
	return t
}

func (t *Tree) List() (l EntryList) {
	return t.Init().list
}

func (t *Tree) String() string {
	return fmt.Sprintf(treeFormat, t.Id(), t.List())
}

func (t *Tree) init() {
	t.list = Ls(t.Id())
}
