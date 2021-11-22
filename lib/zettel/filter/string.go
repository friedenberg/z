package filter

import "github.com/friedenberg/z/lib/zettel"

//TODO remove
type String string

func (f String) FilterZettel(_ int, z *zettel.Zettel) bool {
	if f == "" {
		return true
	} else {
		return z.Metadata.Match(string(f))
	}
}
