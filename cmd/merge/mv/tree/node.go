package tree

import (
	"github.com/emirpasic/gods/trees/avltree"
	"github.com/evovetech/got/cmd/merge/mv/file"
)

type Node avltree.Node

func node(node *avltree.Node) *Node {
	return (*Node)(node)
}

func (n *Node) Entry() Entry {
	return n.Value.(Entry)
}

func (n *Node) Dir() (*Dir, bool) {
	dir, ok := n.Value.(*Dir)
	return dir, ok
}

func (n *Node) File() (*File, bool) {
	f, ok := n.Value.(*File)
	return f, ok
}

func (n *Node) append(parent *Tree, newPath file.Path) (*Tree, int) {
	dir, ok := n.Dir()
	if !ok {
		return nil, -1
	}
	path := dir.key
	i, ok := path.IndexMatch(newPath)
	if !ok {
		return nil, -1
	}

	i++
	if i == len(path) {
		return dir.Tree().PutDir(newPath[i:]), -1
	}

	// Remove node path
	parent.Remove(path)

	// Add new path as base
	oldDir := dir
	dir = NewDir(path[:i])
	parent.Add(dir)
	tree := dir.Tree()

	// update old path, add to base
	oldDir.key = path[i:]
	tree.Add(oldDir)
	return tree, i
}
