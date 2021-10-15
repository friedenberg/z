package pipeline

import (
	"path"
	"strconv"
	"strings"

	"github.com/friedenberg/z/lib"
	"golang.org/x/xerrors"
)

func Import(u lib.Umwelt, oldPath string, shouldCopy bool) (z *lib.Zettel, err error) {
	oldId := strings.TrimSuffix(path.Base(oldPath), path.Ext(oldPath))

	oldIdInt, err := strconv.ParseInt(oldId, 10, 64)

	if err != nil {
		return
	}

	z1 := &lib.Zettel{
		Id:     oldIdInt,
		Path:   oldPath,
		Umwelt: u,
	}

	err = u.Store().Hydrate(z1, true)

	if err != nil {
		return
	}

	ur, hasUrl := z1.Metadata.Url()
	f, hasFile := z1.Metadata.LocalFile()

	if hasUrl && hasFile {
		err = xerrors.Errorf("imported zettel has both url and file")
	} else if hasFile {
		base := path.Dir(oldPath)
		z, err = NewOrFoundForFile(u, f.FilePath(base), shouldCopy)
	} else if hasUrl {
		z, err = NewOrFoundForUrl(u, ur.String())
	} else {
		z, err = New(u)
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
