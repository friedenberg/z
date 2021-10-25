package filter

import "github.com/friedenberg/z/lib"

func True() (f filter) {
	f.filter = func(i int, z *lib.Zettel) bool {
		return true
	}

	return
}
