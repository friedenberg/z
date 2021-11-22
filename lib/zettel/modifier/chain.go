package modifier

import (
	"io"

	"github.com/friedenberg/z/lib/zettel"
	"github.com/friedenberg/z/lib/zettel/beginner"
	"github.com/friedenberg/z/lib/zettel/ender"
)

type chain []Modifier

func (m chain) Begin(w io.Writer) {
	for _, m1 := range m {
		if m2, ok := m1.(beginner.Beginner); ok {
			m2.Begin(w)
		}
	}
}

func (m chain) ModifyZettel(i int, z *zettel.Zettel) (err error) {
	for _, m1 := range m {
		err = m1.ModifyZettel(i, z)

		if err != nil {
			return
		}
	}

	return
}

func (m chain) End(w io.Writer) {
	for _, m1 := range m {
		if m2, ok := m1.(ender.Ender); ok {
			m2.End(w)
		}
	}
}

func Chain(f ...Modifier) (m chain) {
	for _, m1 := range f {
		m = append(m, m1)
	}

	return
}
