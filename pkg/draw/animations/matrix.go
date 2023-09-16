package animations

import (
	"time"
	"trelligo/pkg/draw"
	"trelligo/pkg/shims/rand"
)

var headColor = draw.RGB{200, 255, 200}

const minLength = 3
const maxLength = 8

const colorStep = 250

type scrollLine struct {
	offset int8
	length int8
}
type matrix struct {
	buf         draw.Buffer4x4
	lines       [4]scrollLine
	lastUpdated time.Time
	r           *rand.Rand
}

func NewMatrix(r *rand.Rand) draw.Animation {
	m := &matrix{r: r}

	for i := range &m.lines {
		m.initScrollLine(i)
	}

	return m
}

func (m *matrix) initScrollLine(i int) {
	m.lines[i] = scrollLine{
		offset: int8(-m.r.Intn(16)),
		length: int8(m.r.Intn(maxLength-minLength) + minLength),
	}
}

func (m *matrix) Update(now time.Time) {
	if now.Sub(m.lastUpdated) < time.Millisecond*100 {
		return
	}
	m.lastUpdated = now
	m.buf = draw.Buffer4x4{}

	for column := range &m.lines {
		m.lines[column].offset += 1
		s := m.lines[column]
		if s.offset-s.length > 3 {
			m.initScrollLine(column)
		}
		drawLine(&m.buf, &s, uint8(column))
	}
}

func drawLine(buf *draw.Buffer4x4, s *scrollLine, col uint8) {
	if s.offset-s.length > 3 || s.offset < 0 {
		return
	}

	for i := int8(0); i < s.length; i++ {
		x := s.offset - i
		if !inRange(x) {
			continue
		}
		c := colorAtX(i, s.length)
		buf.Set(uint8(x), col, c)
	}
}

func inRange(n int8) bool {
	return n >= 0 && n < 4
}

func colorAtX(x, length int8) draw.RGB {
	if x == 0 {
		return headColor
	}
	b := uint8(255 - (int(x-1)*colorStep)/int(length))
	return draw.RGB{0, b, 0}
}

func (m *matrix) Draw(d draw.Display) error {
	return d.WriteBuffer(&m.buf)
}
