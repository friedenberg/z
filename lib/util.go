package lib

func WithPrefix(s []string, p string) (a []string) {
	a = make([]string, len(s))
	for i, v := range s {
		a[i] = p + v
	}

	return
}
