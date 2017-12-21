package file

type MoveParser struct {
	from   Path
	to     Path
	result MovePath
}

func NewMove(from Path, to Path) *MoveParser {
	return &MoveParser{
		from: from,
		to:   to,
	}
}

func (p *MoveParser) Parse() (MovePath, bool) {
	for it := p.iterator(); it.hasNext(); {
		p.result.Append(it.get()...)
	}
	return p.result.get()
}

func (p *MoveParser) iterator() *moveIterator {
	return &moveIterator{
		cur: NewPart(p.from.split(), p.to.split()),
	}
}

type moveIterator struct {
	cur  Part
	done bool
}

func (it *moveIterator) hasNext() bool {
	return !it.done
}

func (it *moveIterator) matchEqual() (MovePart, bool) {
	return it.cur.matchEqual()
}

func (it *moveIterator) matchUnequal() (MovePart, bool) {
	return it.cur.matchUnequal()
}

func (it *moveIterator) get() (parts MovePath) {
	done := func() MovePath {
		if last := it.cur; last.Max() > 0 {
			parts.Append(NewMovePart(
				last.From(),
				last.To(),
				false),
			)
		}
		it.cur = nil
		it.done = true
		return parts
	}
	for _, f := range []partMatch{it.matchEqual, it.matchUnequal} {
		if m, ok := f(); ok && parts.Append(m) {
			it.cur = m.next()
		} else {
			return done()
		}
	}
	return
}
