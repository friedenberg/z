package filter

import (
	"github.com/friedenberg/z/lib/zettel"
)

func Not(fs Filter) (f filter) {
	f.filter = func(i int, z *zettel.Zettel) bool {
		return !fs.FilterZettel(i, z)
	}

	return
}
