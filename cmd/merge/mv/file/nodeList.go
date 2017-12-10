package file

//go:generate genny -in=$GOFILE -out=gen-$GOFILE gen "T_=String"

type T_NodeList map[T_]*T_Node

func (list *T_NodeList) init() T_NodeList {
	l := *list
	if l == nil {
		l = make(T_NodeList)
		*list = l
	}
	return l
}

func (list T_NodeList) add(node *T_Node) {
	if n, ok := list[node.Value]; ok {
		for _, child := range node.children {
			n.Add(child)
		}
	} else {
		list[node.Value] = node
	}
}
