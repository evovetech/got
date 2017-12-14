package log

import "io"

type LogWriter struct {
	writer io.Writer
	indent *Indent
}

func NewLogWriter(w io.Writer, i *Indent) *LogWriter {
	return &LogWriter{w, i}
}

func (lw *LogWriter) Write(p []byte) (int, error) {
	if i := lw.indent; i != nil {
		return i.WriteTo(lw.writer, p)
	}
	return lw.writer.Write(p)
}
