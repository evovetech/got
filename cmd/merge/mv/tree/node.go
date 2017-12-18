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

func (n *Node) Dir() (DirEntry, bool) {
	dir, ok := n.Value.(DirEntry)
	return dir, ok
}

func (n *Node) File() (FileEntry, bool) {
	f, ok := n.Value.(FileEntry)
	return f, ok
}

func (n *Node) append(parent DirEntry, newPath file.Path) (DirEntry, int) {
	dir, ok := n.Dir()
	if !ok {
		return nil, -1
	}
	path := dir.Key()
	i, ok := path.IndexMatch(newPath)
	if !ok {
		return nil, -1
	}

	i++
	if i == len(path) {
		return dir.PutDir(newPath[i:]), -1
	}

	// Remove node path
	parent.tree().Remove(path)

	// Add new path as base
	oldDir := dir
	dir = NewDirEntry(path[:i])
	parent.add(dir)

	// update old path, add to base
	oldDir.setKey(path[i:])
	dir.add(oldDir)
	return dir, i
}
