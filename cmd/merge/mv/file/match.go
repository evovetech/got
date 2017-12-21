package file

import "github.com/evovetech/got/util"

func nextMatch(from Path, to Path) (f int, t int, ok bool) {
	if len(from) > len(to) {
		f, t, ok = from.nextMatch(to)
	} else {
		t, f, ok = to.nextMatch(from)
	}
	return
}

func (p Path) nextMatch(o Path) (int, int, bool) {
	oLen := len(o)
	if len(p) == 0 || oLen == 0 {
		return -1, -1, false
	}
	for pi := range p {
		oMax := util.MinInt(pi+1, oLen)
		oi := oMax - 1
		for i := 0; i < pi; i++ {
			if p[i] == o[oi] {
				return i, oi, true
			}
		}
		for oi := 0; oi < oMax; oi++ {
			if p[pi] == o[oi] {
				return pi, oi, true
			}
		}
	}
	return -1, -1, false
}
