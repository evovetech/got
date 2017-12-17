package log

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"sync/atomic"
)

type Indent struct {
	level int32
	size  int32
	buf   []byte
}

func NewIndent(size int32) *Indent {
	indent := new(Indent)
	indent.size = size
	buf := make([]byte, size)
	for i := int32(0); i < size; i++ {
		buf[i] = ' '
	}
	indent.buf = buf
	return indent
}

func (i *Indent) Level() int32 {
	return atomic.LoadInt32(&i.level)
}

func (i *Indent) Size() int32 {
	return i.size
}

func (i *Indent) In() {
	atomic.AddInt32(&i.level, 1)
}

func (i *Indent) Out() {
	for cur := atomic.LoadInt32(&i.level); cur > 0; {
		next := cur - 1
		if atomic.CompareAndSwapInt32(&i.level, cur, next) {
			return
		}
	}
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
		for index := int32(0); index < level; index++ {
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
