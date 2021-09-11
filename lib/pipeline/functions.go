package pipeline

import (
	"github.com/friedenberg/z/lib"
	"github.com/friedenberg/z/util"
	"golang.org/x/xerrors"
)

func NormalizePath(u lib.Umwelt, p string) (n string, err error) {
	n, err = u.FilesAndGit().GetNormalizedPath(p)
	return
}

func HydrateFromIndex(u lib.Umwelt, s string) (z *lib.Zettel, err error) {
	z = &lib.Zettel{Umwelt: u}
	id := util.BaseNameNoSuffix(s)
	zi, ok := u.Index.Get(id)

	if !ok {
		return nil, xerrors.Errorf("missing zettel in index for id '%s'", s)
	}

	u.Index.HydrateZettel(z, zi)

	return
}

func HydrateFromFile(u lib.Umwelt, p string, includeBody bool) (z *lib.Zettel, err error) {
	z = &lib.Zettel{Umwelt: u}
	z.Path = p
	err = z.Hydrate(includeBody)
	return
}
