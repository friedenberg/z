package commands

import (
	"path/filepath"

	"github.com/friedenberg/z/lib"
	"github.com/friedenberg/z/lib/zettel"
	"golang.org/x/xerrors"
)

func DefaultArgNormalizer(u lib.Umwelt) ArgNormalizeFunc {
	return func(_ int, path string) (normalizedArg string, err error) {
		normalizedArg, err = u.Store().GetNormalizedPath(path)
		return
	}
}

func HydrateFromIndexFunc(u lib.Umwelt) HydrateFunc {
	return func(_ int, z *lib.Zettel, path string) error {
		id := filepath.Base(path)

		idId, err := zettel.IdFromString(id)

		if err != nil {
			return err
		}

		zi, ok := u.Index.Get(idId)

		if !ok {
			return xerrors.Errorf("missing zettel in index for id '%s'", path)
		}

		u.Index.HydrateZettel(z, zi)

		return nil
	}
}

func HydrateFromFileFunc(u lib.Umwelt, includeBody bool) HydrateFunc {
	return func(_ int, z *lib.Zettel, path string) error {
		z.Path = path
		return z.Hydrate(includeBody)
	}
}
