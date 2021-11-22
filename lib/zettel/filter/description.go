package filter

import (
	"strings"

	"github.com/friedenberg/z/lib/zettel"
)

type Description string

func (f Description) FilterZettel(_ int, z *zettel.Zettel) bool {
	if f == "" {
		return true
	} else {
		return strings.ToLower(z.Metadata.Description()) == strings.ToLower(string(f))
	}
}
