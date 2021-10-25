package filter

import "github.com/friedenberg/z/lib"

func HasUrl() (f filter) {
	f.filter = func(i int, z *lib.Zettel) bool {
		_, ok := z.Note.Metadata.Url()
		return ok
	}

	return
}
