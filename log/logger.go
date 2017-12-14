package log

import (
	"io"
	"log"
)

type Logger struct {
	*LogWriter
	*log.Logger
}

func New(w io.Writer, prefix string, flags int, indent *Indent) *Logger {
	lw := NewLogWriter(w, indent)
	return &Logger{
		LogWriter: lw,
		Logger:    log.New(lw, prefix, flags),
	}
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
