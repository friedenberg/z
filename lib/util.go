package lib

import (
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

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
	}

	return
}

func TimeFromPath(path string) (t time.Time, err error) {
	base := filepath.Base(path)
	ext := filepath.Ext(path)

	base = base[:len(base)-len(ext)]
	i, err := strconv.ParseInt(base, 10, 64)

	if err != nil {
		return
	}

	t = time.Unix(i, 0)

	return
}
