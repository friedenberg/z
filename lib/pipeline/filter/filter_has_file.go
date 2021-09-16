package filter

import "github.com/friedenberg/z/lib"

func HasFile() (f filter) {
	f.filter = func(i int, z *lib.Zettel) bool {
		return z.Note.Metadata.HasFile()
	}

	return
}
