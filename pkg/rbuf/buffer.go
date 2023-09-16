package rbuf

import (
	"errors"
)

type RingBuffer[T any] struct {
	buf    []T
	iwrite int
	iread  int
}

var ErrBufferOverflow = errors.New("buffer is full")
var ErrBufferUnderflow = errors.New("buffer is empty")

func New[T any](size int) RingBuffer[T] {
	return RingBuffer[T]{
		buf: make([]T, size),
	}
}
func (r *RingBuffer[T]) Write(t T) error {
	i := (r.iwrite + 1) % len(r.buf)
	if i == r.iread {
		return ErrBufferOverflow
	}
	r.buf[i] = t
	r.iwrite = i
	return nil
}

func (r *RingBuffer[T]) Read() (T, error) {
	if r.iread == r.iwrite {
		var t T
		return t, ErrBufferUnderflow
	}
	r.iread = (r.iread + 1) % len(r.buf)
	return r.buf[r.iread], nil
}
