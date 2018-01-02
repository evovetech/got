package log

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/evovetech/got/util"
	"io"
)

type Indent struct {
	level util.Counter
	size  uint32
	buf   []byte
}

func NewIndent(size uint32) *Indent {
	indent := &Indent{
		level: util.NewCounter(),
		size:  size,
	}
	buf := make([]byte, size)
	for i := uint32(0); i < size; i++ {
		buf[i] = ' '
	}
	indent.buf = buf
	return indent
}

func (i *Indent) Level() uint32 {
	return i.level.Get()
}

func (i *Indent) Size() uint32 {
	return i.size
}

func (i *Indent) In() {
	i.level.IncrementAndGet()
}

func (i *Indent) Out() {
	i.level.DecrementAndGet()
}

func (i *Indent) Transform(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	if _, err := i.WriteTo(&buf, data); err != nil {
		return data, err
	}
	return buf.Bytes(), nil
}

func (i *Indent) WriteTo(w io.Writer, data []byte) (int, error) {
	s := bufio.NewScanner(bytes.NewReader(data))
	for s.Scan() {
		level := i.Level()
		for index := uint32(0); index < level; index++ {
			if n, err := w.Write(i.buf); err != nil {
				return n, err
			}
		}
		if n, err := w.Write(s.Bytes()); err != nil {
			return n, err
		}
		if n, err := fmt.Fprintln(w, ""); err != nil {
			return n, err
		}
	}
	return len(data), nil
}
