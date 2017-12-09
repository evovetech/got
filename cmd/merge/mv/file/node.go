package file

import (
	"github.com/cheekybits/genny/generic"
)

//go:generate genny -in=$GOFILE -out=gen-$GOFILE gen "T_=String"

type T_ generic.Type

type T_Node struct {
	Value    T_
	Children T_NodeList
}
