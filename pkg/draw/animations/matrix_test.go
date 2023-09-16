package animations

import (
	"fmt"
	"strings"
	"testing"
	"time"
	"trelligo/pkg/be"
	"trelligo/pkg/draw"
	"trelligo/pkg/shims/rand"
)

type mockDisplay struct {
	buffers []*draw.Buffer4x4
}

func (m *mockDisplay) String() string {

	var sbuf strings.Builder
	sbuf.WriteString("=======================\n")
	for _, buf := range m.buffers {
		sbuf.WriteString(formatBuffer(buf))
		sbuf.WriteString("=======================\n")
	}
	return sbuf.String()
}

func (m *mockDisplay) WriteBuffer(b *draw.Buffer4x4) error {

	fmt.Println(formatBuffer(b))

	var buf draw.Buffer4x4
	copy(buf[:], b[:])
	m.buffers = append(m.buffers, &buf)
	return nil
}

func formatBuffer(buf *draw.Buffer4x4) string {

	var sbuf strings.Builder
	for y := 0; y < 4; y++ {
		for x := 0; x < 4; x++ {
			c := buf[x*4+y]
			if c == headColor {
				sbuf.WriteString("H ")
			} else if c.G > 0 {
				sbuf.WriteString("G ")
			} else {
				sbuf.WriteString(". ")
			}
		}
		sbuf.WriteString("\n")
	}
	return sbuf.String()
}

func TestMatrix_Draw(t *testing.T) {

	r := rand.New(rand.NewSource(1))

	fmt.Printf("starting\n")
	m := NewMatrix(r)

	d := &mockDisplay{}

	start := time.Now()

	for i := 0; i < 24; i++ {
		fmt.Printf("update\n")
		m.Update(start)
		fmt.Printf("draw\n")
		err := m.Draw(d)
		be.NoError(t, err)

		start = start.Add(time.Millisecond * 100)
	}

	fmt.Println(d.String())
}
