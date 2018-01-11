package commit

import (
	"github.com/evovetech/got/git"
	"github.com/evovetech/got/git/object"
	"github.com/evovetech/got/git/tree"
	"regexp"
)

var (
	reCommitLine = regexp.MustCompile("^(\\w+)\\s+(.*)$")
)

type Commit struct {
	*object.Object

	tree    *tree.Tree
	parents List
}

func New(id object.Id) *Commit {
	c := &Commit{Object: object.NewCommit(id)}
	c.SetInitFunc(c.init)
	return c
}

func (c *Commit) Init() *Commit {
	c.Object.Init()
	return c
}

func (c *Commit) Tree() *tree.Tree {
	return c.Init().tree
}

func (c *Commit) Parents() *List {
	return &c.Init().parents
}

func (c *Commit) init() {
	git.Command("cat-file", "-p", c.Id().String()).ForEachLine(func(line string) error {
		if match := reCommitLine.FindStringSubmatch(line); match != nil {
			switch match[1] {
			case "tree":
				c.tree = tree.New(object.Id(match[2]))
			case "parent":
				p := New(object.Id(match[2]))
				c.parents.Append(p)
			}
		}
		return nil
	})
}
