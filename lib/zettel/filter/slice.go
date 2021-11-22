package filter

type Filters []Filter

func (f Filters) Or() (f1 Or) {
	f1 = Or(f)
	return
}
