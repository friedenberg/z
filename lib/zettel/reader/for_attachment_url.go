package reader

import (
	"net/url"

	"github.com/friedenberg/z/lib"
	"github.com/friedenberg/z/lib/zettel"
	"github.com/friedenberg/z/lib/zettel/metadata"
)

func ForAttachmentUrl() reader {
	return Make(
		func(u *lib.Umwelt, i int, b []byte) (*zettel.Zettel, error) {
			return newOrFoundForUrl(u, i, string(b))
		},
	)
}

func newOrFoundForUrl(u *lib.Umwelt, i int, urlString string) (z *zettel.Zettel, err error) {
	ur, err := url.Parse(urlString)

	if err != nil {
		return
	}

	id, ok := u.Index.Urls.GetId(urlString)

	if ok {
		z, err = hydrateFromFile(u, id.String()+".md", true)
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
