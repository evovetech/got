package tree

import (
	"github.com/evovetech/got/file"
	"github.com/evovetech/got/git"
	"github.com/evovetech/got/git/types"
	"github.com/evovetech/got/log"
	"github.com/evovetech/got/util"
	"regexp"
)

var (
	reTreeLine = regexp.MustCompile("^(\\d+) (\\w+) (\\w+)\\t(.*)$")
)

type Entry struct {
	Mode string
	Kind types.Type
	Sha  types.Sha
	Path file.Path
}

func ParseTree(sha types.Sha) file.Dir {
	dir := file.NewRoot()
	git.Command("ls-tree", sha.String()).ForEachLine(func(line string) error {
		if match := reTreeLine.FindStringSubmatch(line); match != nil {
			entry := newEntry(match)
			log.Println(util.String(entry))
		}
		return nil
	})
	return dir
}

func newEntry(match []string) *Entry {
	return &Entry{
		Mode: match[1],
		Kind: types.Parse(match[2]),
		Sha:  types.Sha(match[3]),
		Path: file.GetPath(match[4]),
	}
}
