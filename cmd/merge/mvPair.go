package merge

import (
	"path"

	"github.com/evovetech/got/git"
)

type MvPair struct {
	From string
	To   string
}

func (mv *MvPair) run() error {
	if err := mkdir(path.Dir(mv.To)); err != nil {
		return err
	}
	return git.Command("mv", mv.From, mv.To).Run()
}
