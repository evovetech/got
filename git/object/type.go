package object

import (
	"encoding/json"
	"strings"
)

type (
	Type uint8
)

const (
	None Type = iota
	Commit
	Tree
	Blob
	Tag
)

var typeStrings = []string{
	"",       /* OBJ_NONE = 0 */
	"commit", /* OBJ_COMMIT = 1 */
	"tree",   /* OBJ_TREE = 2 */
	"blob",   /* OBJ_BLOB = 3 */
	"tag",    /* OBJ_TAG = 4 */
}

func Parse(t string) Type {
	t = strings.ToLower(t)
	for i, s := range typeStrings {
		if s == t {
			return Type(i)
		}
	}
	return None
}

func (t Type) String() string {
	return typeStrings[t]
}

func (t Type) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.String())
}
