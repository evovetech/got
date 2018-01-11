package object

import (
	"github.com/evovetech/got/git"
	"github.com/evovetech/got/git/types"
	"regexp"
)

var (
	reCommitLine = regexp.MustCompile("^(\\w+)\\s+(.*)$")
)

type commit struct {
	Object

	tree    Tree
	parents CommitList
}

func NewCommit(id types.Id) Commit {
	c := &commit{Object: New(id, types.Commit)}
	c.SetInitFunc(c.populate)
	return c
}

func (c *commit) Tree() Tree {
	return c.init().tree
}

func (c *commit) Parents() *CommitList {
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
				c.tree = NewTree(types.Id(match[2]))
			case "parent":
				p := NewCommit(types.Id(match[2]))
				c.parents.Append(p)
			}
		}
		return nil
	})
}
