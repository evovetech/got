package mv

import (
	"github.com/evovetech/got/util"
)

type Group struct {
	FileName string
	From     []FilePath
	To       []FilePath
}

func (g *Group) String() string {
	return util.String(g)
}

func (g *Group) IsValid() bool {
	return len(g.From) > 0 && len(g.To) > 0
}

func (g *Group) Add(fp FilePath) {
	g.To = append(g.To, fp)
}

func (g *Group) Del(fp FilePath) {
	g.From = append(g.From, fp)
}

func (g *Group) isSize(size int) bool {
	return len(g.From) == size && len(g.To) == size
}

func (g *Group) isSameSize() bool {
	return len(g.From) == len(g.To)
}

func (g *Group) first() Rename {
	return Rename{
		From: g.From[0],
		To:   g.To[0],
	}
}

func (g *Group) parse(firstTry bool) (*Group, []Rename) {
	if g.isSize(1) {
		return nil, []Rename{
			g.first(),
		}
	}
	if firstTry {
		var mvPairs []Rename
		var leftFrom = g.From
		var leftTo []FilePath
		for _, to := range g.To {
			var index = -1
			var from FilePath
			toDir := to.Dir()
			for i, fr := range leftFrom {
				if toDir.MovedFrom(fr.Dir()) {
					index = i
					from = fr
					break
				}
			}
			if index == -1 {
				leftTo = append(leftTo, to)
			} else {
				mvPair := Rename{
					From: from,
					To:   to,
				}
				mvPairs = append(mvPairs, mvPair)
				leftFrom = append(leftFrom[:index], leftFrom[index+1:]...)
			}
		}
		var err *Group
		if len(leftTo) > 0 {
			var mvPair []Rename
			err, mvPair = (&Group{
				FileName: g.FileName,
				From:     leftFrom,
				To:       leftTo,
			}).parse(false)
			if len(mvPair) > 0 {
				mvPairs = append(mvPairs, mvPair...)
			}
		}
		return err, mvPairs
	}
	return g, nil
}
