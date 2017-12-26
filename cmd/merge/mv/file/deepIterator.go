package file

type DeepIterator interface {
	Iterator
	Parent() Path
	Dir() Path
	FullPath() Path
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
	return state.it
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
	return d.state.Key()
}

func (d *deepIterator) Value() interface{} {
	return d.state.Entry()
}

func (d *deepIterator) Path() Path {
	return d.state.Path()
}

func (d *deepIterator) Entry() Entry {
	return d.state.Entry()
}

func (d *deepIterator) Dir() Path {
	return d.state.Dir()
}

func (d *deepIterator) Parent() Path {
	return d.state.Parent()
}

func (d *deepIterator) FullPath() Path {
	if state := d.state; state != nil {
		return JoinPaths(state.Parent(), state.Dir(), state.Path())
	}
	return nil
}

type itState struct {
	parent *itState
	dir    Dir
	it     Iterator

	cur Entry
}

func rootState(dir Dir) *itState {
	return &itState{dir: dir, it: dir.Iterator()}
}

func (s *itState) Key() interface{} {
	if s == nil {
		return nil
	}
	return s.it.Key()
}

func (s *itState) Entry() Entry {
	if s == nil {
		return nil
	}
	return s.cur
}

func (s *itState) Path() Path {
	if e := s.Entry(); e != nil {
		return e.Path()
	}
	return nil
}

func (s *itState) Dir() Path {
	if s == nil {
		return nil
	}
	return s.dir.Path()
}

func (s *itState) Parent() Path {
	if s == nil {
		return nil
	}
	if p := s.parent; p != nil {
		return JoinPaths(p.Parent(), p.Dir())
	}
	return nil
}

func (s *itState) next() (*itState, bool) {
	if s == nil {
		return nil, false
	}
	if state, found := s.nextLevel(); found {
		state.it.Begin()
		return state.next()
	}
	if s.it.Next() {
		s.cur = s.it.Entry()
		return s, true
	}
	s.cur = nil
	return s.parent.next()
}

func (s *itState) prev() (*itState, bool) {
	if s == nil {
		return nil, false
	}
	if state, found := s.nextLevel(); found {
		state.it.End()
		return state.prev()
	}
	if s.it.Prev() {
		s.cur = s.it.Entry()
		return s, true
	}
	s.cur = nil
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
		state = &itState{parent: s, dir: dir, it: dir.Iterator()}
	}
	return
}
