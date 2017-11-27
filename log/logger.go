package log

import (
	"io"
	"log"
)

type Logger struct {
	*log.Logger
}

func New(w io.Writer, prefix string, flags int) *Logger {
	return &Logger{
		Logger: log.New(w, prefix, flags),
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
