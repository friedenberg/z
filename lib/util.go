package lib

import (
	"fmt"
	"path/filepath"
	"strconv"
	"time"
)

func TimeFromPath(path string) (t time.Time, err error) {
	base := filepath.Base(path)
	ext := filepath.Ext(path)

	base = base[:len(base)-len(ext)]
	i, err := strconv.ParseInt(base, 10, 64)

	if err != nil {
		err = fmt.Errorf("time from path: %w", err)
		return
	}

	t = time.Unix(i, 0)

	return
}

func ZettelIdFromPath(path string) (zi string, err error) {
	t, err := TimeFromPath(path)

	if err != nil {
		err = fmt.Errorf("zettel id from path: %w", err)
		return
	}

	i := t.Unix()

	zi = strconv.FormatInt(i, 10)

	return
}
