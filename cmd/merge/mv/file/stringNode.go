package file

import (
	"path/filepath"
	"strings"
)

func ParseString(path string) *StringNode {
	slashy := filepath.ToSlash(path)
	paths := strings.Split(slashy, "/")
	return ParseStringPath(paths)
}
