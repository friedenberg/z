package filter

import (
	"github.com/friedenberg/z/lib"
	"github.com/friedenberg/z/lib/pipeline"
)

func Not(fs pipeline.Filter) (f filter) {
	f.filter = func(i int, z *lib.Zettel) bool {
		return !fs.FilterZettel(i, z)
	}

	return
}
