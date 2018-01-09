package object

type Tree interface {
	Object
}

type Commit interface {
	Object

	Tree() Tree
	Parents() CommitList
}

type CommitList interface {
	Item() Commit
	Next() CommitList
}
