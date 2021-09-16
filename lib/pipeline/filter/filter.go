package filter

import "github.com/friedenberg/z/lib"

type FilterFunc func(int, *lib.Zettel) bool

type filter struct {
	filter FilterFunc
}

func (f filter) FilterZettel(i int, z *lib.Zettel) bool {
	return f.filter(i, z)
}

func Make(ff FilterFunc) (f filter) {
	f.filter = ff
	return
}
