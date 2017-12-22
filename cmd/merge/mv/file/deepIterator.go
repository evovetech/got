package file

import "container/list"

type DeepIterator interface {
	Iterator
}

type deepIterator struct {
	root  Iterator
	stack itStack
}

func (d *deepIterator) reset() Iterator {
	d.stack.reset(d.root)
	return d.root
}

func (d *deepIterator) Begin() {
	d.reset().Begin()
}

func (d *deepIterator) End() {
	d.reset().End()
}

func (d *deepIterator) First() bool {
	return d.reset().First()
}

func (d *deepIterator) Last() bool {
	return d.reset().Last()
}

func (d *deepIterator) Next() bool {
	stack := d.stack
	it := stack.it()
	switch e := it.Entry().(type) {
	case Dir:
		stack.push(e.Iterator())
	default:
		if it.Next() {
			return true
		}
		stack.pop()
	}
	return d.Next()
}

func (d *deepIterator) Prev() bool {
	stack := d.stack
	it := stack.it()
	switch e := it.Entry().(type) {
	case Dir:
		it := e.Iterator()
		it.End()
		stack.push(it)
	default:
		if it.Prev() {
			return true
		}
		stack.pop()
	}
	return d.Prev()
}

func (d *deepIterator) Key() interface{} {
	return d.stack.it().Key()
}

func (d *deepIterator) Value() interface{} {
	return d.stack.it().Value()
}

func (d *deepIterator) Path() Path {
	return d.stack.it().Path()
}

func (d *deepIterator) Entry() Entry {
	return d.stack.it().Entry()
}

type itStack struct {
	stack list.List
}

func (s *itStack) reset(it Iterator) {
	s.stack.Init()
	s.stack.PushFront(it)
}

func (s *itStack) cur() (e *list.Element, it Iterator) {
	if e = s.stack.Front(); e != nil {
		it, _ = e.Value.(Iterator)
	}
	if it == nil {
		it = noEntries
	}
	return
}

func (s *itStack) it() Iterator {
	_, it := s.cur()
	return it
}

func (s *itStack) push(it Iterator) {
	s.stack.PushFront(it)
}

func (s *itStack) pop() Iterator {
	if e, it := s.cur(); e != nil {
		s.stack.Remove(e)
		return it
	}
	return noEntries
}
