package filter

import (
	"strings"

	"github.com/friedenberg/z/lib/zettel"
)

type tag string

func Tag(q string) (t tag) {
	t = tag(strings.ToLower(q))
	return
}

func (f tag) FilterZettel(_ int, z *zettel.Zettel) bool {
	if f == "" {
		return true
	}

	_, ok := z.Metadata.SearchMatchTags().Get(string(f))
	return ok
}

//TODO-P3 add set method to filter.Tag
// func (f *tag) Set(t string) (err error) {
// 	f = Tag(t)
// 	return
// }
