package log

import (
	"log"
	"os"
)

var (
	indent  = NewIndent(2)
	Verbose = New(new(verbose), "", 0, indent)
	Std     = New(os.Stdout, "", 0, indent)
	Err     = New(os.Stderr, "", log.Llongfile, indent)
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

func IndentIn() int {
	return indent.In()
}

func IndentOut() int {
	return indent.Out()
}
