package file

import (
	"encoding/json"
	"github.com/cheekybits/genny/generic"
)

//go:generate genny -in=$GOFILE -out=gen-$GOFILE gen "T_=String"

type T_ generic.Type

type T_List = []T_

type T_Node struct {
	Value T_

	children T_NodeList
}

func NewT_Node(val T_) *T_Node {
	return &T_Node{Value: val}
}

func ParseT_Path(path T_List) *T_Node {
	var node *T_Node
	if l := len(path); l > 0 {
		node = NewT_Node(path[0])
		if l > 1 {
			child := ParseT_Path(path[1:])
			node.Add(child)
		}
	}
	return node
}

var RootT_Node T_Node

func (n *T_Node) getChildren() T_NodeList {
	return n.children.init()
}

func (n *T_Node) Add(node *T_Node) {
	n.getChildren().add(node)
}

func (n T_Node) MarshalJSON() ([]byte, error) {
	return json.Marshal(n.children)
}
