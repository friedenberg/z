package filter

import "github.com/friedenberg/z/lib/zettel"

func True() (f filter) {
	f.filter = func(i int, z *zettel.Zettel) bool {
		return true
	}

	return
}
