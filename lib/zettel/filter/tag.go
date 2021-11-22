package filter

import "github.com/friedenberg/z/lib/zettel"

type Tag string

func (f Tag) FilterZettel(_ int, z *zettel.Zettel) bool {
	if f == "" {
		return true
	} else {
		return z.Metadata.Match(string(f))
	}
}
