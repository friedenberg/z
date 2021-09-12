package commands

import (
	"path/filepath"

	"github.com/friedenberg/z/lib"
	"github.com/friedenberg/z/lib/pipeline"
	"golang.org/x/xerrors"
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

func MatchQuery(q string) pipeline.Filter {
	return func(i int, z *lib.Zettel) bool {
		return doesZettelMatchQuery(z, q)
	}
}
