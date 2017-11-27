package log

import (
	"log"
	"os"
)

var (
	Verbose = New(new(verbose), "", 0)
	Std     = New(os.Stdout, "", 0)
	Err     = New(os.Stderr, "", log.Llongfile)
)

func Print(v interface{}) {
	Verbose.Print(v)
}

func Println(v interface{}) {
	Verbose.Println(v)
}

func Printf(format string, v ...interface{}) {
	Verbose.Printf(format, v...)
}
