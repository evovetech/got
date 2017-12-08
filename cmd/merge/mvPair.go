package merge

import (
	"github.com/evovetech/got/git"
)

type MvPair struct {
	From FilePath
	To   FilePath
}

func (mv *MvPair) run() error {
	if err := mv.To.Dir().Mkdirs(); err != nil {
		return err
	}
	return git.Command("mv", mv.From.String(), mv.To.String()).Run()
}
