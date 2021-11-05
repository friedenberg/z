package modifier

import (
	"io"

	"github.com/friedenberg/z/lib/pipeline"
	"github.com/friedenberg/z/lib/zettel"
)

type chain []pipeline.Modifier

func (m chain) Begin(w io.Writer) {
	for _, m1 := range m {
		if m2, ok := m1.(pipeline.Beginner); ok {
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
		if m2, ok := m1.(pipeline.Ender); ok {
			m2.End(w)
		}
	}
}

func Chain(f ...pipeline.Modifier) (m chain) {
	for _, m1 := range f {
		m = append(m, m1)
	}

	return
}
