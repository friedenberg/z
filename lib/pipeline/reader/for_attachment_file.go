package reader

import (
	"github.com/friedenberg/z/lib"
	"github.com/friedenberg/z/lib/zettel/metadata"
)

func ForAttachmentFile() reader {
	return Make(
		func(u lib.Umwelt, i int, b []byte) (*lib.Zettel, error) {
			return newForFile(u, i, string(b))
		},
	)
}

func newForFile(u lib.Umwelt, i int, file string) (z *lib.Zettel, err error) {
	z, err = readerNew(u, i, file)

	if err != nil {
		return
	}

	fd := metadata.NewFile{
		Path: file,
	}

	z.Note.Metadata.AddStringTags(fd.Tag())

	return
}
