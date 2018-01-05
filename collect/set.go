package collect

import (
	"github.com/cheekybits/genny/generic"
)

//go:generate genny -in=$GOFILE -out=gen-$GOFILE gen "T_=String,Sha"

type T_ generic.Type

type T_Set map[T_]bool

func (set *T_Set) Init() (s T_Set) {
	if s = *set; s == nil {
		s = make(T_Set)
		*set = s
	}
	return
}

func (set *T_Set) Add(val T_) bool {
	if s := set.Init(); !s[val] {
		s[val] = true
		return true
	}
	return false
}
