package reader

import (
	"path"
	"strconv"
	"strings"

	"github.com/friedenberg/z/lib"
	"github.com/friedenberg/z/lib/zettel"
	"golang.org/x/xerrors"
)

func Import(shouldCopy bool) reader {
	return Make(
		func(u *lib.Umwelt, i int, b []byte) (*zettel.Zettel, error) {
			return importZettel(u, i, string(b), shouldCopy)
		},
	)
}

func importZettel(u *lib.Umwelt, i int, oldPath string, shouldCopy bool) (z *zettel.Zettel, err error) {
	oldId := strings.TrimSuffix(path.Base(oldPath), path.Ext(oldPath))

	oldIdInt, err := strconv.ParseInt(oldId, 10, 64)

	if err != nil {
		return
	}

	z1 := &zettel.Zettel{
		Id:      zettel.Id(oldIdInt),
		Path:    oldPath,
		ZUmwelt: u,
	}

	err = u.Kasten.Hydrate(u, z1, true)

	if err != nil {
		return
	}

	ur, hasUrl := z1.Metadata.Url()
	f, hasFile := z1.Metadata.LocalFile()

	if hasUrl && hasFile {
		err = xerrors.Errorf("imported zettel has both url and file")
	} else if hasFile {
		base := path.Dir(oldPath)
		z, err = newForFile(u, i, f.FilePath(base))
	} else if hasUrl {
		z, err = newOrFoundForUrl(u, i, ur.String())
	} else {
		z, err = readerNew(u, i, "")
	}

	if err != nil {
		return
	}

	z.Merge(z1)

	if !shouldCopy {
		//TODO-P3 should delete imported zettel?
	}

	return
}
