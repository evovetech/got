package file

//go:generate genny -in=$GOFILE -out=gen-$GOFILE gen "T_=String,Int"

type T_NodeList map[T_]*T_Node

func (list T_NodeList) Add(node *T_Node) {
	if n, ok := list[node.Value]; ok {
		for _, child := range node.Children {
			n.Children.Add(child)
		}
	} else {
		list[node.Value] = node
	}
}
