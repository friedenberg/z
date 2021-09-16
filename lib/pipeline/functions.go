package pipeline

import (
	"net/url"

	"github.com/friedenberg/z/lib"
	"golang.org/x/xerrors"
)

func NormalizePath(u lib.Umwelt, p string) (n string, err error) {
	n, err = u.FilesAndGit().GetNormalizedPath(p)
	return
}

func HydrateFromIndex(u lib.Umwelt, s string) (z *lib.Zettel, err error) {
	z = &lib.Zettel{Umwelt: u}

	id, err := lib.IdFromString(s)

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
	z = &lib.Zettel{Umwelt: u}
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

	if ok && len(ids) > 1 {
		err = xerrors.Errorf("multiple zettels ('%q') with url: '%s'", ids, urlString)
		return
	} else if ok && len(ids) == 1 {
		z, err = HydrateFromIndex(u, ids[0].String())
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

//func NewOrFoundForFile(u lib.Umwelt, file string) (z *lib.Zettel, err error) {
//	//TODO check if file exists on disk
//	//TODO check if file sha exists in cache
//	_, err = url.Parse(urlString)

//	if err != nil {
//		return
//	}

//	ids, ok := u.Index.Urls.Get(urlString, u.Index)

//	if ok && len(ids) > 1 {
//		err = xerrors.Errorf("multiple zettels ('%q') with url: '%s'", ids, urlString)
//		return
//	} else if ok && len(ids) == 1 {
//		z, err = HydrateFromIndex(u, ids[0].String())
//		return
//	}

//	z, err = New(u)

//	if err != nil {
//		return
//	}

//	//TODO check if Metadata exists
//	z.Metadata.Url = urlString

//	return
//}

func New(u lib.Umwelt) (z *lib.Zettel, err error) {
	id, err := u.FilesAndGit().NewId()

	if err != nil {
		return
	}

	z = &lib.Zettel{
		Umwelt: u,
		Id:     id.Int(),
		Path:   lib.MakePathFromId(u.FilesAndGit().BasePath, id.String()),
	}

	return
}

// func New(u lib.Umwelt) (z *lib.Zettel, err error) {
// 	id, err := u.FilesAndGit().NewId()

// 	if err != nil {
// 		return
// 	}

// 	z = &lib.Zettel{
// 		Umwelt: u,
// 		Id:     id.Int(),
// 		Path:   lib.MakePathFromId(u.FilesAndGit().BasePath, id.String()),
// 	}

// 	return
// }
