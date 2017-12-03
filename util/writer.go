package util

import "io"

type compositeWriter struct {
	writers []io.Writer
}

func (w *compositeWriter) Write(b []byte) (n int, err error) {
	main := w.writers[0]
	if n, err = main.Write(b); err == nil {
		for _, w := range w.writers[1:] {
			w.Write(b)
		}
	}
	return
}

func CompositeWriter(writers ...io.Writer) io.Writer {
	var ret []io.Writer
	for _, w := range writers {
		if w != nil {
			ret = append(ret, w)
		}
	}
	switch len(ret) {
	case 0:
		return nil
	case 1:
		return ret[0]
	default:
		return &compositeWriter{ret}
	}
}
