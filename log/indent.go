package log

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"sync"
)

type Indent struct {
	level int
	buf   []byte
	size  []byte
	rw    sync.RWMutex
}

func NewIndent(size int) *Indent {
	i := new(Indent)
	if size < 0 {
		size = 0
	}
	i.size = make([]byte, size)
	for index := range i.size {
		i.size[index] = ' '
	}
	return i
}

func (i *Indent) Prefix() string {
	i.rw.RLock()
	defer i.rw.RUnlock()
	return string(i.buf)
}

func (i *Indent) In() int {
	i.rw.Lock()
	defer i.rw.Unlock()
	i.level++
	return i.updateBuf()
}

func (i *Indent) Out() int {
	i.rw.Lock()
	defer i.rw.Unlock()
	i.level--
	return i.updateBuf()
}

func (i *Indent) Transform(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	if _, err := i.WriteTo(&buf, data); err != nil {
		return data, err
	}
	return buf.Bytes(), nil
}

func (i *Indent) WriteTo(w io.Writer, data []byte) (num int, err error) {
	i.rw.RLock()
	defer i.rw.RUnlock()
	s := bufio.NewScanner(bytes.NewReader(data))
	for s.Scan() {
		if _, err = w.Write(i.buf); err != nil {
			return
		}
		var n int
		if n, err = w.Write(s.Bytes()); err != nil {
			return
		}
		num += n
		if n, err = fmt.Fprintln(w); err != nil {
			return
		}
		num += n
	}
	return
}

func (i *Indent) updateBuf() int {
	size := i.level * len(i.size)
	for cap(i.buf) < size {
		i.buf = append(i.buf, i.size...)
	}
	i.buf = i.buf[:size]
	return size
}
