package file

type DeepIterator interface {
	Iterator
	Dir() Path
}

type deepIterator struct {
	root  Dir
	state *itState
}

func newDeepIterator(dir Dir) DeepIterator {
	d := &deepIterator{root: dir}
	d.reset()
	return d
}

func (d *deepIterator) reset() Iterator {
	state := rootState(d.root)
	d.state = state
	return state.it()
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
	if next, found := d.state.next(); found {
		d.state = next
		return true
	}
	d.state = nil
	return false
}

func (d *deepIterator) Prev() bool {
	if prev, found := d.state.prev(); found {
		d.state = prev
		return true
	}
	d.state = nil
	return false
}

func (d *deepIterator) Key() interface{} {
	return d.state.it().Key()
}

func (d *deepIterator) Value() interface{} {
	return d.state.Entry()
}

func (d *deepIterator) Dir() Path {
	return d.state.Dir()
}

func (d *deepIterator) Path() Path {
	return d.state.Path()
}

func (d *deepIterator) Entry() Entry {
	return d.state.Entry()
}

type itState struct {
	parent *itState
	dir    Dir
	itr    Iterator

	cur Entry
}

func rootState(dir Dir) *itState {
	return &itState{dir: dir, itr: dir.Iterator()}
}

func (s *itState) Dir() Path {
	if s == nil {
		return nil
	}
	return s.dir.Path()
}

func (s *itState) Entry() Entry {
	if s == nil {
		return nil
	}
	return s.cur
}

func (s *itState) Path() Path {
	if e := s.Entry(); e != nil {
		// TODO: parent path
		if dir := s.Dir(); dir != nil {
			return dir.Append(e.Path())
		}
		return e.Path()
	}
	return nil
}

func (s *itState) it() Iterator {
	if s == nil {
		return noEntries
	}
	if it := s.itr; it != nil {
		return it
	}
	return noEntries
}

func (s *itState) next() (*itState, bool) {
	if s == nil {
		return nil, false
	}
	if state, found := s.nextLevel(); found {
		state.itr.Begin()
		return state.next()
	}
	if s.itr.Next() {
		s.cur = s.itr.Entry()
		return s, true
	}
	return s.parent.next()
}

func (s *itState) prev() (*itState, bool) {
	if s == nil {
		return nil, false
	}
	if state, found := s.nextLevel(); found {
		state.itr.End()
		return state.prev()
	}
	if s.itr.Prev() {
		s.cur = s.itr.Entry()
		return s, true
	}
	return s.parent.prev()
}

func (s *itState) nextLevel() (state *itState, found bool) {
	if s.cur == nil {
		return nil, false
	}
	cur := s.cur
	s.cur = nil
	switch dir := cur.(type) {
	case Dir:
		found = true
		state = &itState{parent: s, dir: dir, itr: dir.Iterator()}
	}
	return
}
