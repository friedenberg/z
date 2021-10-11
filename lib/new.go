package lib

import (
	"path"
	"strconv"
	"time"
)

func MakePathFromId(basePath, id string) string {
	return path.Join(basePath, id+".md")
}

func MakePathFromTime(basePath string, t time.Time) (filename string) {
	unixTime := t.Unix()
	return MakePathFromId(basePath, strconv.FormatInt(unixTime, 10))
}
