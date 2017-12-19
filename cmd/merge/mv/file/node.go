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

func (n *Node) Dir() (Dir, bool) {
	dir, ok := n.Value.(Dir)
	return dir, ok
}

func (n *Node) File() (File, bool) {
	f, ok := n.Value.(File)
	return f, ok
}

func (n *Node) match(path Path) (dir Dir, offset int, ok bool) {
	if dir, ok = n.Dir(); ok {
		dirPath := dir.Key()
		if offset, ok = dirPath.IndexMatch(path); ok {
			offset++
			offset -= len(dirPath)
		}
	}
	return
}

func (n *Node) find(child Path) (Entry, bool) {
	if parent, offset, ok := n.match(child); ok {
		if path := parent.Key(); offset == 0 {
			i := len(path) + offset
			if e, ok := parent.Find(path[:i]); ok {
				return e, ok
			}
		}
	}
	return nil, false
}

func (n *Node) append(parent Dir, child Path) (Dir, int) {
	dir, offset, ok := n.match(child)
	if !ok {
		return nil, -1
	}

	path := dir.Key()
	i := offset + len(path)
	if offset == 0 {
		return dir.PutDir(child[i:]), -1
	}

	// Remove node path
	parent.tree().Remove(path)

	// Add new path as base
	oldDir := dir
	dir = NewDir(path[:i])
	parent.put(dir)

	// update old path, add to base
	oldDir.setPath(path[i:])
	dir.put(oldDir)
	return dir, i
}
