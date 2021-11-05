package filter

import "github.com/friedenberg/z/lib/zettel"

func HasUrl() (f filter) {
	f.filter = func(i int, z *zettel.Zettel) bool {
		_, ok := z.Note.Metadata.Url()
		return ok
	}

	return
}
