package filter

import (
	"github.com/friedenberg/z/lib/zettel"
)

func MatchQueries(qs ...string) (f filter) {
	f.filter = func(i int, z *zettel.Zettel) bool {
		for _, q := range qs {
			if q == "" {
				continue
			}

			if z.Metadata.Match(q) {
				return true
			}
		}

		return false
	}

	return
}
