package object

import (
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

func NewCommit(id Id) Commit {
	c := &commit{Object: New(id, CommitType)}
	c.SetInitFunc(func() {
		c.tree, c.parents, _ = catCommit(id)
	})
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
