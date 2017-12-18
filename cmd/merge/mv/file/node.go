package file

import (
	"github.com/emirpasic/gods/trees/avltree"
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

func (n *Node) File() (File, bool) {
	f, ok := n.Value.(File)
	return f, ok
}

func (n *Node) append(parent DirEntry, newPath Path) (DirEntry, int) {
	dir, ok := n.Dir()
	if !ok {
		return nil, -1
	}
	path := dir.Path()
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
	oldDir.setPath(path[i:])
	dir.add(oldDir)
	return dir, i
}
