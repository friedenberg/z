package reader

import (
	"github.com/friedenberg/z/lib"
	"github.com/friedenberg/z/lib/zettel"
	"golang.org/x/xerrors"
)

func FromIndex(u lib.Umwelt, _ int, s string) (z *zettel.Zettel, err error) {
	z = &zettel.Zettel{
		ZUmwelt: u,
	}

	id, err := zettel.IdFromString(s)

	if err != nil {
		return
	}

	zi, ok := u.Index.Get(id)

	if !ok {
		z = nil
		err = xerrors.Errorf("missing zettel in index for id '%s'", s)
		return
	}

	u.Index.HydrateZettel(z, zi)

	return
}
