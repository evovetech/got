package object

import (
	"github.com/evovetech/got/git"
	"github.com/evovetech/got/types"
)

type Id = types.Sha

type Object interface {
	Id() Id
	Type() git.Type
	String() string
	MarshalJSON() ([]byte, error)

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

	List() List
}
