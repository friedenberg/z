package filter

import "github.com/friedenberg/z/lib"

func MatchQueries(qs ...string) (f filter) {
	f.filter = func(i int, z *lib.Zettel) bool {
		for _, q := range qs {
			if q == "" {
				continue
			}

			if !z.Metadata.TagSet().Match(q) {
				return false
			}
		}

		return true
	}

	return
}
