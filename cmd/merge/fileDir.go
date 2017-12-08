package merge

import (
	"path/filepath"
	"strings"
)

type FileDir string

func (d FileDir) Base() string {
	return filepath.Base(string(d))
}

func (d FileDir) MovedFrom(other FileDir) bool {
	dStr := string(d)
	oStr := string(other)
	// TODO:
	return d.Base() == other.Base() ||
		strings.Contains(dStr, other.Base()) ||
		strings.Contains(oStr, d.Base())
}
