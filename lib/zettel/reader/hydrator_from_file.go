package reader

import (
	"github.com/friedenberg/z/lib"
	"github.com/friedenberg/z/lib/zettel"
)

func FromFile(includeBody bool) reader {
	return Make(
		func(u *lib.Umwelt, _ int, b []byte) (*zettel.Zettel, error) {
			return hydrateFromFile(u, string(b), includeBody)
		},
	)
}

func hydrateFromFile(u *lib.Umwelt, p string, includeBody bool) (z *zettel.Zettel, err error) {
	z = &zettel.Zettel{
		ZUmwelt: u,
		Path:    p,
	}

	err = u.Kasten.Hydrate(u, z, includeBody)

	return
}
