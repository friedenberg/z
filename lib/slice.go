package lib

import "github.com/friedenberg/z/lib/zettel"

type ZettelSlice []*zettel.Zettel

func (s ZettelSlice) Paths() (p []string) {
	p = make([]string, len(s))

	for i, z := range s {
		p[i] = z.Path
	}

	return
}
