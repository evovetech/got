package log

import "io"

type Writer struct {
	io.Writer
	*Indent
}

func NewWriter(w io.Writer, i *Indent) *Writer {
	return &Writer{w, i}
}

func (lw *Writer) Write(p []byte) (int, error) {
	if i := lw.Indent; i != nil {
		return i.WriteTo(lw.Writer, p)
	}
	return lw.Writer.Write(p)
}
