package mv

import (
	"github.com/evovetech/got/cmd/merge/mv/file"
	"github.com/evovetech/got/git"
)

type Rename struct {
	From file.Path
	To   file.Path
}

func (p *Rename) run() error {
	if err := p.To.Dir().Mkdirs(); err != nil {
		return err
	}
	return git.Command("mv", p.From.String(), p.To.String()).Run()
}
