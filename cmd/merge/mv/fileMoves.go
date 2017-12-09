package mv

import (
	"github.com/evovetech/got/git"
	"github.com/evovetech/got/log"
	"github.com/evovetech/got/util"
	"regexp"
)

var reAdd = regexp.MustCompile("^A\\s+(.*)")
var reDel = regexp.MustCompile("^D\\s+(.*)")
var reRename = regexp.MustCompile("^R\\s+(.*)\\s+->\\s+(.*)")

type FileMoves struct {
	Renames []Rename
	errs    []*Group
}

func GetFileMoves() (*FileMoves, bool) {
	errs, renames := NewMap().Run()
	if len(renames) == 0 {
		if len(errs) != 0 {
			log.Verbose.Printf("errors: %s", util.String(errs))
		}
		return nil, true
	}
	return &FileMoves{
		Renames: renames,
		errs:    errs,
	}, false
}

func (m *FileMoves) Run() error {
	// abort merge, move files
	git.AbortMerge()
	for _, mv := range m.Renames {
		mv.run()
	}
	if len(m.errs) > 0 {
		log.Verbose.Printf("errors: %s", util.String(m.errs))
	}
	return git.Command("commit", "-m", "moving files to prepare for merge").Run()
}
