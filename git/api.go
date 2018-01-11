package git

import (
	"github.com/evovetech/got/git/types"
)

type Object interface {
	Id() types.Id
	Type() types.Type
	String() string

	SetInitFunc(func())
	Init()
}

type Commit interface {
	Object

	Tree() Tree
	Parents() *CommitList
}

type Tree interface {
	Object

	List() ObjectList
}
