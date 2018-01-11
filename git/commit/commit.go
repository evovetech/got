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

type Commit interface {
	object.Object

	Tree() tree.Tree
	Parents() *List
}

type commit struct {
	object.Object

	tree    tree.Tree
	parents List
}

func New(id object.Id) Commit {
	c := &commit{Object: object.NewCommit(id)}
	c.SetInitFunc(c.populate)
	return c
}

func (c *commit) Tree() tree.Tree {
	return c.init().tree
}

func (c *commit) Parents() *List {
	return &c.init().parents
}

func (c *commit) init() *commit {
	c.Object.Init()
	return c
}

func (c *commit) populate() {
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
