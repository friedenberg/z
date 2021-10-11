package pipeline

import (
	"net/url"
	"os/exec"
	"path"
	"strconv"
	"strings"

	"github.com/friedenberg/z/lib"
	"github.com/friedenberg/z/lib/zettel"
	"github.com/friedenberg/z/lib/zettel/metadata"
	"github.com/friedenberg/z/util"
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
		Path:   p,
	}

	err = z.Hydrate(includeBody)

	return
}

func NewOrFoundForUrl(u lib.Umwelt, urlString string) (z *lib.Zettel, err error) {
	ur, err := url.Parse(urlString)

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

	urlTag := metadata.Url(*ur)
	z.Metadata.SetUrl(urlTag)

	return
}

func NewOrFoundForFile(u lib.Umwelt, file string, shouldCopy bool) (z *lib.Zettel, err error) {
	sum, err := util.Sha256HashForFile(file)

	if err != nil {
		return
	}

	ids, ok := u.Index.Files.Get(sum, u.Index)

	if ok && ids.Len() > 1 {
		err = xerrors.Errorf("multiple zettels ('%q') with file: '%s'", ids, sum)
		return
	} else if ok && ids.Len() == 1 {
		z, err = HydrateFromIndex(u, ids.Slice()[0].String())
		return
	}

	z, err = New(u)

	if err != nil {
		return
	}

	//TODO-P0 check for checksum file name collisions
	n := sum[0:7]

	fd := metadata.File{
		Id:  n,
		Ext: strings.ReplaceAll(util.ExtNoDot(file), "-", ""),
	}

	if shouldCopy {
		cmd := exec.Command("cp", "-R", file, fd.FilePath(u.BasePath))
		msg, err := cmd.CombinedOutput()

		if err != nil {
			err = xerrors.Errorf("%w: %s", err, msg)
		}
	} else {
		cmd := exec.Command("mv", file, fd.FilePath(u.BasePath))
		msg, err := cmd.CombinedOutput()

		if err != nil {
			err = xerrors.Errorf("%w: %s", err, msg)
		}
	}

	if err != nil {
		return
	}

	z.Note.Metadata.AddFile(fd)

	return
}

func Import(u lib.Umwelt, oldPath string, shouldCopy bool) (z *lib.Zettel, err error) {
	oldId := strings.TrimSuffix(path.Base(oldPath), path.Ext(oldPath))

	oldIdInt, err := strconv.ParseInt(oldId, 10, 64)

	if err != nil {
		return
	}

	z1 := &lib.Zettel{
		Id:     oldIdInt,
		Path:   oldPath,
		Umwelt: &u,
	}

	err = z1.Hydrate(true)

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
		Note: lib.Note{
			Metadata: metadata.MakeMetadata(),
		},
	}

	z.Metadata.AddStringTags(u.TagsForNewZettels...)

	return
}
