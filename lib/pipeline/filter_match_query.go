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

		//TODO-P1 normalize
		if _, ok := z.Metadata.TagSet().Get(q); ok {
			return true
		}

		return false
	}
}
