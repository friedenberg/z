package pipeline

import "github.com/friedenberg/z/lib"

func And(a, b Filter) Filter {
	return func(i int, z *lib.Zettel) bool {
		return a(i, z) && b(i, z)
	}
}
