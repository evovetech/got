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
	Std.Print(v)
}

func Println(v interface{}) {
	Std.Println(v)
}

func Printf(format string, v ...interface{}) {
	Std.Printf(format, v...)
}
