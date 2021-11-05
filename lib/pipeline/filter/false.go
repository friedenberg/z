package filter

import "github.com/friedenberg/z/lib/zettel"

func False() (f filter) {
	f.filter = func(i int, z *zettel.Zettel) bool {
		return false
	}

	return
}
