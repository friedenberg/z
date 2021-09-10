package commands

import (
	"fmt"
	"path/filepath"

	"github.com/friedenberg/z/lib"
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
			return fmt.Errorf("missing zettel in index for id '%s'", path)
		}

		z.Id = zi.Id
		z.Path = zi.Path
		z.Metadata = zi.Metadata

		return nil
	}
}

func HydrateFromFileFunc(u lib.Umwelt) HydrateFunc {
	return func(_ int, z *lib.Zettel, path string) error {
		z.Path = path
		return z.Hydrate(true)
	}
}
