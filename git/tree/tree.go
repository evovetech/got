package tree

import (
	"fmt"
	"github.com/evovetech/got/git/types"
	"regexp"
)

var (
	reTreeLine = regexp.MustCompile("^(\\d+) (\\w+) (\\w+)\\t(.*)$")
)

const treeFormat = `{
  sha: "%s",
  list: %s,
}
`

type Tree struct {
	sha  types.Sha
	list EntryList
}

func NewTree(sha string) *Tree {
	return &Tree{sha: types.Sha(sha)}
}

func (t Tree) Sha() types.Sha {
	return t.sha
}

func (t Tree) String() string {
	return fmt.Sprintf(treeFormat, t.Sha(), t.List())
}

func (t Tree) List() (l EntryList) {
	if l = t.list; l == nil {
		l = Ls(t.sha)
		t.list = l
	}
	return
}
