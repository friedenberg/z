package alfred

import "strings"

func WithPrefix(s []string, p string) (a []string) {
	a = make([]string, len(s))
	for i, v := range s {
		a[i] = p + v
	}

	return
}

func Split(s []string, sep string) (a []string) {
	a = make([]string, len(s))

	for _, v := range s {
		for _, b := range strings.Split(v, sep) {
			a = append(a, b)
		}

		a = append(a, v)
	}

	return
}
