package merge

import (
	"github.com/evovetech/got/git"
	"github.com/evovetech/got/util"
)

// TODO:
type mergeRef struct {
	Original git.Ref
	Temp     git.Ref
}

func (m *mergeRef) Update() error {
	return m.Temp.Update()
}

func (m *mergeRef) String() string {
	return util.String(m)
}
