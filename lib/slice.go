package lib

type ZettelSlice []*Zettel

func (s ZettelSlice) Paths() (p []string) {
	p = make([]string, len(s))

	for i, z := range s {
		p[i] = z.Path
	}

	return
}
