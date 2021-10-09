package pipeline

import (
	"net/url"
	"os/exec"
	"path"
	"strconv"

	"github.com/friedenberg/z/lib"
	"github.com/friedenberg/z/lib/zettel"
	"golang.org/x/xerrors"
)

func NormalizePath(u lib.Umwelt, p string) (n string, err error) {
	n, err = u.FilesAndGit().GetNormalizedPath(p)
	return
}

func HydrateFromIndex(u lib.Umwelt, s string) (z *lib.Zettel, err error) {
	z = &lib.Zettel{
		Umwelt: &u,
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
		Umwelt: &u,
	}
	z.Path = p
	err = z.Hydrate(includeBody)
	return
}

func NewOrFoundForUrl(u lib.Umwelt, urlString string) (z *lib.Zettel, err error) {
	_, err = url.Parse(urlString)

	if err != nil {
		return
	}

	ids, ok := u.Index.Urls.Get(urlString, u.Index)

	if ok && ids.Len() > 1 {
		err = xerrors.Errorf("multiple zettels ('%q') with url: '%s'", ids, urlString)
		return
	} else if ok && ids.Len() == 1 {
		z, err = HydrateFromIndex(u, ids.Slice()[0].String())
		return
	}

	z, err = New(u)

	if err != nil {
		return
	}

	//TODO check if Metadata exists
	z.Metadata.Url = urlString

	return
}

func NewOrFoundForFile(u lib.Umwelt, file string, shouldCopy bool) (z *lib.Zettel, err error) {
	//TODO check if file exists on disk
	//TODO check if file sha exists in cache
	z, err = New(u)

	if err != nil {
		return
	}

	fileName := strconv.FormatInt(z.Id, 10) + path.Ext(file)
	fileName, err = NormalizePath(u, fileName)

	if err != nil {
		return
	}

	if shouldCopy {
		cmd := exec.Command("cp", "-R", file, fileName)
		msg, err := cmd.CombinedOutput()

		if err != nil {
			err = xerrors.Errorf("%w: %s", err, msg)
		}
	} else {
		cmd := exec.Command("mv", file, fileName)
		msg, err := cmd.CombinedOutput()

		if err != nil {
			err = xerrors.Errorf("%w: %s", err, msg)
		}
	}

	if err != nil {
		return
	}

	//TODO check if Metadata exists
	z.Metadata.File = fileName

	return
}

func New(u lib.Umwelt) (z *lib.Zettel, err error) {
	id, err := u.FilesAndGit().NewId()

	if err != nil {
		return
	}

	z = &lib.Zettel{
		Id:     id.Int(),
		Path:   lib.MakePathFromId(u.FilesAndGit().BasePath, id.String()),
		Umwelt: &u,
	}

	return
}
