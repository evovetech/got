package file

import (
	"github.com/emirpasic/gods/containers"
	"github.com/emirpasic/gods/trees/avltree"
)

type Iterator interface {
	containers.ReverseIteratorWithKey
	Path() Path
	Entry() Entry
}

type emptyIterator uint8

const noEntries = emptyIterator(0)

func (emptyIterator) Prev() bool         { return false }
func (emptyIterator) End()               {}
func (emptyIterator) Last() bool         { return false }
func (emptyIterator) Next() bool         { return false }
func (emptyIterator) Value() interface{} { return nil }
func (emptyIterator) Key() interface{}   { return nil }
func (emptyIterator) Begin()             {}
func (emptyIterator) First() bool        { return false }
func (emptyIterator) Path() Path         { return nil }
func (emptyIterator) Entry() Entry       { return nil }

type iterator struct {
	containers.ReverseIteratorWithKey
}

func newIterator(tree *avltree.Tree) Iterator {
	return &iterator{tree.Iterator()}
}

func (it *iterator) Path() Path {
	if p, ok := it.Key().(Path); ok {
		return p
	}
	return nil
}

func (it *iterator) Entry() Entry {
	if e, ok := it.Value().(Entry); ok {
		return e
	}
	return nil
}
