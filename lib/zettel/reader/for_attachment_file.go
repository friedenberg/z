package reader

import (
	"github.com/friedenberg/z/lib"
	"github.com/friedenberg/z/lib/zettel"
	"github.com/friedenberg/z/lib/zettel/metadata"
)

func ForAttachmentFile() reader {
	return Make(
		func(u lib.Umwelt, i int, b []byte) (*zettel.Zettel, error) {
			return newForFile(u, i, string(b))
		},
	)
}

func newForFile(u lib.Umwelt, i int, file string) (z *zettel.Zettel, err error) {
	z, err = readerNew(u, i, file)

	if err != nil {
		return
	}

	fd := metadata.NewFile{
		Path: file,
	}

	z.Metadata.SetFile(&fd)

	return
}
