package commands

import (
	"path/filepath"

	"github.com/friedenberg/z/lib"
	"golang.org/x/xerrors"
)

var (
	HydrateFromIndex HydrateFunc
)

func DefaultArgNormalizer(u lib.Umwelt) ArgNormalizeFunc {
	return func(_ int, path string) (normalizedArg string, err error) {
		normalizedArg, err = u.FilesAndGit().GetNormalizedPath(path)
		return
	}
}

func HydrateFromIndexFunc(u lib.Umwelt) HydrateFunc {
	return func(_ int, z *lib.Zettel, path string) error {
		id := filepath.Base(path)
		zi, ok := u.Index.Get(id)

		if !ok {
			return xerrors.Errorf("missing zettel in index for id '%s'", path)
		}

		z.Id = zi.Id
		z.Path = zi.Path
		z.Metadata = zi.Metadata

		return nil
	}
}

func HydrateFromFileFunc(u lib.Umwelt, includeBody bool) HydrateFunc {
	return func(_ int, z *lib.Zettel, path string) error {
		z.Path = path
		return z.Hydrate(includeBody)
	}
}
