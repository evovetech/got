package tree

import (
	"github.com/emirpasic/gods/utils"
	"github.com/evovetech/got/cmd/merge/mv/file"
	"github.com/evovetech/got/util"
)

func PathComparator(a, b interface{}) int {
	ap := a.(file.Path)
	bp := b.(file.Path)
	min := util.MinInt(len(ap), len(bp))
	comp := utils.StringComparator
	for i := 0; i < min; i++ {
		if diff := comp(ap[i], bp[i]); diff != 0 {
			return diff
		}
	}
	return len(ap) - len(bp)
}
