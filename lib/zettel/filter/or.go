package filter

import (
	"github.com/friedenberg/z/lib/pipeline"
	"github.com/friedenberg/z/lib/zettel"
)

func Or(fs ...pipeline.Filter) (f filter) {
	f.filter = func(i int, z *zettel.Zettel) bool {
		for _, f1 := range fs {
			if f1 == nil {
				continue
			}

			if f1.FilterZettel(i, z) {
				return true
			}
		}

		return false
	}

	return
}
