package file

import (
	"github.com/cheekybits/genny/generic"
)

//go:generate genny -in=$GOFILE -out=gen-$GOFILE gen "T_=String,Int"

type T_ generic.Type

type T_Path = []T_

type T_Node struct {
	Value    T_
	Children T_NodeList
}

var RootT_Node T_Node

func init() {
	RootT_Node.init()
}

func (n *T_Node) init() *T_Node {
	if n.Children == nil {
		n.Children = make(T_NodeList)
	}
	return n
}

func NewT_Node(val T_) *T_Node {
	n := &T_Node{Value: val}
	return n.init()
}

func ParseT_Node(path T_Path) *T_Node {
	var node *T_Node
	if l := len(path); l > 0 {
		node = NewT_Node(path[0])
		if l > 1 {
			child := ParseT_Node(path[1:])
			node.Children.Add(child)
		}
	}
	return node
}
