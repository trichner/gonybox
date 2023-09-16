package rbuf

import (
	"testing"
	"trelligo/pkg/be"
)

func TestRingBuffer_Read(t *testing.T) {
	capacity := 12
	buf := RingBuffer[int]{
		buf: make([]int, capacity),
	}

	_, err := buf.Read()
	be.Equal(t, err, ErrBufferUnderflow)

	v := 10
	err = buf.Write(v)
	be.NoError(t, err)

	a, err := buf.Read()
	be.NoError(t, err)
	be.Equal(t, a, v)
}

func TestRingBuffer_Cap(t *testing.T) {
	size := 12
	buf := RingBuffer[int]{
		buf: make([]int, size),
	}

	// actual capacity is actually one less because we did the lazy implementation
	for i := 0; i < 11; i++ {
		err := buf.Write(i)
		be.NoError(t, err)
	}

	err := buf.Write(-1)
	be.Equal(t, err, ErrBufferOverflow)
}
