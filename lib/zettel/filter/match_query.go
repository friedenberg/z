package filter

//TODO-P4 refactor
func MatchQuery(q string) (f filter) {
	if q == "" {
		f = True()
	} else {
		f = MatchQueries(q)
	}

	return
}
