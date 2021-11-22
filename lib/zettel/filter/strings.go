package filter

type Strings []string

func (f Strings) Filters() (fs Filters) {
	fs = make([]Filter, len(f))

	for i, f := range f {
		fs[i] = String(f)
	}

	return
}
