package commit

import (
	"github.com/evovetech/got/git"
	"github.com/evovetech/got/git/object"
	"github.com/evovetech/got/git/tree"
	"github.com/evovetech/got/git/types"
	"regexp"
)

var (
	reCommitLine = regexp.MustCompile("^(\\w+)\\s+(.*)$")
)

type commit struct {
	git.Object

	tree    git.Tree
	parents git.CommitList
}

func New(id types.Id) git.Commit {
	c := &commit{Object: object.NewCommit(id)}
	c.SetInitFunc(c.populate)
	return c
}

func (c *commit) Tree() git.Tree {
	return c.init().tree
}

func (c *commit) Parents() *git.CommitList {
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
				c.tree = tree.New(types.Id(match[2]))
			case "parent":
				p := New(types.Id(match[2]))
				c.parents.Append(p)
			}
		}
		return nil
	})
}
