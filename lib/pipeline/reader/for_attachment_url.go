package reader

import (
	"net/url"

	"github.com/friedenberg/z/lib"
	"github.com/friedenberg/z/lib/zettel/metadata"
	"golang.org/x/xerrors"
)

func ForAttachmentUrl() reader {
	return Make(
		func(u lib.Umwelt, i int, b []byte) (*lib.Zettel, error) {
			return newOrFoundForUrl(u, i, string(b))
		},
	)
}

func newOrFoundForUrl(u lib.Umwelt, i int, urlString string) (z *lib.Zettel, err error) {
	ur, err := url.Parse(urlString)

	if err != nil {
		return
	}

	ids, ok := u.Index.Urls.GetIds(urlString, u.Index)

	if ok && ids.Len() > 1 {
		err = xerrors.Errorf("multiple zettels ('%q') with url: '%s'", ids, urlString)
		return
	} else if ok && ids.Len() == 1 {
		z, err = hydrateFromFile(u, ids.Slice()[0].String()+".md", true)
		return
	}

	z, err = readerNew(u, i, "")

	if err != nil {
		return
	}

	urlTag := metadata.Url{URL: *ur}
	z.Metadata.SetUrl(urlTag)

	return
}
