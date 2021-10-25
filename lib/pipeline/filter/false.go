package filter

import "github.com/friedenberg/z/lib"

func False() (f filter) {
	f.filter = func(i int, z *lib.Zettel) bool {
		return false
	}

	return
}
