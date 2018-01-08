package collect

//go:generate genny -in=$GOFILE -out=gen-$GOFILE gen "T_=String,Sha"

type T_List []T_

func (list *T_List) Append(t T_) {
	*list = append(*list, t)
}
