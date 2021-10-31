package filter

import "github.com/friedenberg/z/lib"

//TODO-P4 refactor
func MatchQuery(q string) (f filter) {
	f.filter = func(i int, z *lib.Zettel) bool {
		if q == "" {
			return true
		}

		//TODO-P2
		// if z.Note.Metadata.LocalFile() == q {
		// 	return true
		// }

		// if z.Metadata.Url == q {
		// 	return true
		// }

		return z.Metadata.TagSet().Match(q)
	}

	return
}