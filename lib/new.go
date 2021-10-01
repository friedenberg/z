package lib

import (
	"path"
	"strconv"
	"time"
)

//TODO deprecate
func (z *Zettel) InitFromTime(t time.Time) {
	z.Path = MakePathFromTime(z.Kasten.Local.BasePath, t)
	z.Id = t.Unix()

	z.Metadata = Metadata{
		Date: t.Format("2006-01-02"),
	}

	return
}

func MakePathFromId(basePath, id string) string {
	return path.Join(basePath, id+".md")
}

func MakePathFromTime(basePath string, t time.Time) (filename string) {
	unixTime := t.Unix()
	return MakePathFromId(basePath, strconv.FormatInt(unixTime, 10))
}
