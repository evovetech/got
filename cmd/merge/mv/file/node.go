package file

import (
	"github.com/cheekybits/genny/generic"
	"encoding/json"
)

//go:generate genny -in=$GOFILE -out=gen-$GOFILE gen "T_=String,Int"

type T_ generic.Type

type T_Path = []T_

type T_Node struct {
	Value    T_
	Children T_NodeList
}

func NewT_Node(val T_) *T_Node {
	n := &T_Node{Value: val}
	return n.init()
}

func ParseT_Path(path T_Path) *T_Node {
	var node *T_Node
	if l := len(path); l > 0 {
		node = NewT_Node(path[0])
		if l > 1 {
			child := ParseT_Path(path[1:])
			node.Children.Add(child)
		}
	}
	return node
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

func (n *T_Node) Add(node *T_Node) {
	n.Children.Add(node)
}

func (n *T_Node) MarshalJSON() ([]byte, error) {
	return json.Marshal(n.Children)
}
