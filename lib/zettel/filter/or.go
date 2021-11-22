package filter

import (
	"github.com/friedenberg/z/lib/zettel"
)

type Or []Filter

func MakeOr(fs ...Filter) (f Or) {
	f = Or(fs)

	return
}

func (f Or) FilterZettel(i int, z *zettel.Zettel) bool {
	for _, f1 := range f {
		if f1 == nil {
			continue
		}

		if f1.FilterZettel(i, z) {
			return true
		}
	}

	return false
}
