package filter

import "github.com/friedenberg/z/lib/zettel"

func HasFile() (f filter) {
	f.filter = func(i int, z *zettel.Zettel) bool {
		return z.Note.Metadata.HasFile()
	}

	return
}
