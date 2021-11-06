package filter

import "github.com/friedenberg/z/lib/zettel"

type Filter interface {
	FilterZettel(int, *zettel.Zettel) bool
}

type FilterFunc func(int, *zettel.Zettel) bool

type filter struct {
	filter FilterFunc
}

func (f filter) FilterZettel(i int, z *zettel.Zettel) bool {
	return f.filter(i, z)
}

func Make(ff FilterFunc) (f filter) {
	f.filter = ff
	return
}
