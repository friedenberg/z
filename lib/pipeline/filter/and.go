package filter

import (
	"github.com/friedenberg/z/lib"
	"github.com/friedenberg/z/lib/pipeline"
)

func And(fs ...pipeline.Filter) (f filter) {
	f.filter = func(i int, z *lib.Zettel) bool {
		for _, f1 := range fs {
			if f1 == nil {
				continue
			}

			if !f1.FilterZettel(i, z) {
				return false
			}
		}

		return true
	}

	return
}
