package object

import (
	"encoding/json"
	"strings"
)

type (
	Type uint8
)

const (
	NoneType Type = iota
	CommitType
	TreeType
	BlobType
	TagType
)

var typeStrings = []string{
	"",       /* OBJ_NONE = 0 */
	"commit", /* OBJ_COMMIT = 1 */
	"tree",   /* OBJ_TREE = 2 */
	"blob",   /* OBJ_BLOB = 3 */
	"tag",    /* OBJ_TAG = 4 */
}

func ParseType(t string) Type {
	t = strings.ToLower(t)
	for i, s := range typeStrings {
		if s == t {
			return Type(i)
		}
	}
	return NoneType
}

func (t Type) Creator() Creator {
	switch t {
	case CommitType:
		return newCommit
	case TreeType:
		return newTree
	case BlobType:
		return newBlob
	case TagType:
		// TODO:
	}
	return newObject
}

func (t Type) New(id Id) Object {
	return t.Creator()(id)
}

func (t Type) String() string {
	return typeStrings[t]
}

func (t Type) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.String())
}

func newObject(id Id) Object {
	return New(id, NoneType)
}

func newCommit(id Id) Object {
	return NewCommit(id)
}

func newTree(id Id) Object {
	return NewTree(id)
}

func newBlob(id Id) Object {
	return NewBlob(id)
}
