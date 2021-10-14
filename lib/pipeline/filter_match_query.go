package pipeline

import "github.com/friedenberg/z/lib"

//TODO-P4 refactor
func MatchQuery(q string) Filter {
	return func(i int, z *lib.Zettel) bool {
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

		for _, t := range z.Note.Metadata.TagStrings() {
			if t == q {
				return true
			}
		}

		return false
	}
}
