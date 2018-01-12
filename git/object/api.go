package object

import (
	"github.com/evovetech/got/types"
)

type Id = types.Sha

type Object interface {
	Id() Id
	Type() Type

	Init()
	String() string
	MarshalJSON() ([]byte, error)

	// private
	setInitFunc(func())
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

type Blob interface {
	Object

	Contents() ([]byte, error)
}
