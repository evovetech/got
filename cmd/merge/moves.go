package merge

import (
	"github.com/evovetech/got/git"
	"github.com/evovetech/got/log"
	"github.com/evovetech/got/util"
	"regexp"
)

var reAdd = regexp.MustCompile("^A\\s+(.*)")
var reDel = regexp.MustCompile("^D\\s+(.*)")
var reRename = regexp.MustCompile("^R\\s+(.*)\\s+->\\s+(.*)")

type Moves struct {
	Renames []MvPair
	errs    []*AddDel
}

func getMoves() (*Moves, bool) {
	var renames []MvPair
	var adm = make(AddDelMap)
	for _, status := range git.Command("status", "-s", "--untracked-files=all").OutputLines() {
		switch {
		case reAdd.MatchString(status):
			match := reAdd.FindStringSubmatch(status)
			adm.Do(match, Add)
		case reDel.MatchString(status):
			match := reDel.FindStringSubmatch(status)
			adm.Do(match, Del)
		case reRename.MatchString(status):
			match := reRename.FindStringSubmatch(status)
			from := match[1]
			to := match[2]
			mv := MvPair{From: from, To: to}
			renames = append(renames, mv)
			log.Printf("Rename: %s", util.String(mv))
		}
	}
	errs, mvs := adm.Parse()
	renames = append(renames, mvs...)

	if len(renames) == 0 {
		if len(errs) != 0 {
			log.Printf("errors: %s", util.String(errs))
		}
		return nil, true
	}
	return &Moves{
		Renames: renames,
		errs:    errs,
	}, false
}

func (m *Moves) Run() error {
	// abort merge, move files
	git.AbortMerge()
	for _, mv := range m.Renames {
		mv.run()
	}
	log.Printf("errors: %s", util.String(m.errs))
	return git.Command("commit", "-m", "moving files to prepare for merge").Run()
}
