package pipeline

import (
	"net/url"

	"github.com/friedenberg/z/lib"
	"github.com/friedenberg/z/lib/zettel/metadata"
	"golang.org/x/xerrors"
)

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
