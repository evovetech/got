package log

import (
	"bytes"
	"io"
	"log"
)

type Logger struct {
	*Writer
	*log.Logger
}

func NewLogger(w *Writer, prefix string, flags int) *Logger {
	return &Logger{
		Writer: w,
		Logger: log.New(w, prefix, flags),
	}
}

func NewBufLogger(buffer *bytes.Buffer) *Logger {
	return NewLogger(NewWriter(buffer, NewIndent(2)), "", 0)
}

func New(w io.Writer, prefix string, flags int, indent *Indent) *Logger {
	return NewLogger(NewWriter(w, indent), prefix, flags)
}

func (l *Logger) Enter(prefix interface{}, f func(*Logger)) {
	l.Printf("%s {\n", prefix)
	l.In()
	f(l)
	l.Out()
	l.Println("}")
}

func (l *Logger) Write(p []byte) (int, error) {
	var str string
	if str = string(p); str != "" {
		if err := l.Output(3, str); err != nil {
			return 0, err
		}
	}
	return len(p), nil
}
