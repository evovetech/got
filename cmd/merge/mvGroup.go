package merge

import (
	"github.com/evovetech/got/util"
)

type MvGroup struct {
	FileName string
	From     []FilePath
	To       []FilePath
}

func (mv *MvGroup) String() string {
	return util.String(mv)
}

func (mv *MvGroup) IsValid() bool {
	return len(mv.From) > 0 && len(mv.To) > 0
}

func (mv *MvGroup) Add(fp FilePath) {
	mv.To = append(mv.To, fp)
}

func (mv *MvGroup) Del(fp FilePath) {
	mv.From = append(mv.From, fp)
}

func (mv *MvGroup) isSize(size int) bool {
	return len(mv.From) == size && len(mv.To) == size
}

func (mv *MvGroup) isSameSize() bool {
	return len(mv.From) == len(mv.To)
}

func (mv *MvGroup) first() MvPair {
	return MvPair{
		From: mv.From[0],
		To:   mv.To[0],
	}
}

func (mv *MvGroup) parse(firstTry bool) (*MvGroup, []MvPair) {
	if mv.isSize(1) {
		return nil, []MvPair{
			mv.first(),
		}
	}
	if firstTry {
		var mvPairs []MvPair
		var leftFrom = mv.From
		var leftTo []FilePath
		for _, to := range mv.To {
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
				mvPair := MvPair{
					From: from,
					To:   to,
				}
				mvPairs = append(mvPairs, mvPair)
				leftFrom = append(leftFrom[:index], leftFrom[index+1:]...)
			}
		}
		var err *MvGroup
		if len(leftTo) > 0 {
			var mvPair []MvPair
			err, mvPair = (&MvGroup{
				FileName: mv.FileName,
				From:     leftFrom,
				To:       leftTo,
			}).parse(false)
			if len(mvPair) > 0 {
				mvPairs = append(mvPairs, mvPair...)
			}
		}
		return err, mvPairs
	}

	return mv, nil
}
