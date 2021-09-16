package reader

import (
	"github.com/friedenberg/z/lib"
)

func FromFile(includeBody bool) reader {
	return Make(
		func(u lib.Umwelt, _ int, b []byte) (*lib.Zettel, error) {
			return hydrateFromFile(u, string(b), includeBody)
		},
	)
}

func hydrateFromFile(u lib.Umwelt, p string, includeBody bool) (z *lib.Zettel, err error) {
	z = &lib.Zettel{
		Umwelt: u,
		Path:   p,
	}

	err = u.Kasten.Hydrate(z, includeBody)

	return
}
