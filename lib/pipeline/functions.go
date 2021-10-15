package pipeline

import (
	"github.com/friedenberg/z/lib"
	"github.com/friedenberg/z/lib/zettel"
	"golang.org/x/xerrors"
)

func NormalizePath(u lib.Umwelt, p string) (n string, err error) {
	n, err = u.Store().GetNormalizedPath(p)
	return
}

func HydrateFromIndex(u lib.Umwelt, s string) (z *lib.Zettel, err error) {
	z = &lib.Zettel{
		Umwelt: u,
	}

	id, err := zettel.IdFromString(s)

	if err != nil {
		return
	}

	zi, ok := u.Index.Get(id)

	if !ok {
		return nil, xerrors.Errorf("missing zettel in index for id '%s'", s)
	}

	u.Index.HydrateZettel(z, zi)

	return
}

func HydrateFromFile(u lib.Umwelt, p string, includeBody bool) (z *lib.Zettel, err error) {
	z = &lib.Zettel{
		Umwelt: u,
		Path:   p,
	}

	err = u.Store().Hydrate(z, includeBody)

	return
}
