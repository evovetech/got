package log

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"sync/atomic"
)

type Indent struct {
	level Counter
	size  uint32
	buf   []byte
}

type Counter interface {
	Get() uint32
	IncrementAndGet() uint32
	GetAndIncrement() uint32
	DecrementAndGet() uint32
	TryDecrementAndGet() (uint32, bool)
}

func NewCounter() Counter {
	return new(counter)
}

type counter uint32

func (c *counter) ptr() *uint32 {
	return (*uint32)(c)
}

func (c *counter) Get() uint32 {
	return atomic.LoadUint32(c.ptr())
}

func (c *counter) IncrementAndGet() uint32 {
	return c.GetAndIncrement() + 1
}

func (c *counter) GetAndIncrement() uint32 {
	if cur := c.Get(); c.CompareAndSwap(cur, cur+1) {
		return cur
	}
	return c.GetAndIncrement()
}

func (c *counter) DecrementAndGet() uint32 {
	if v, ok := c.TryDecrementAndGet(); ok {
		return v
	}
	return c.DecrementAndGet()
}

func (c *counter) TryDecrementAndGet() (uint32, bool) {
	if cur := c.Get(); cur > 0 {
		next := cur + ^uint32(0)
		return next, c.CompareAndSwap(cur, next)
	}
	return 0, false
}

func (c *counter) CompareAndSwap(cur uint32, next uint32) bool {
	return atomic.CompareAndSwapUint32(c.ptr(), cur, next)
}

func NewIndent(size uint32) *Indent {
	indent := &Indent{
		level: NewCounter(),
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
