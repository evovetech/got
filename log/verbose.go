package log

import (
	"os"

	"github.com/evovetech/got/options"
)

type verbose int8

func (l *verbose) Write(p []byte) (n int, err error) {
	if !options.Verbose {
		return len(p), nil
	}
	return os.Stdout.Write(p)
}
