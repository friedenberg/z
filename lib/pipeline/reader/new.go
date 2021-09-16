package reader

import (
	"github.com/friedenberg/z/lib"
	"github.com/friedenberg/z/lib/zettel/metadata"
)

func New() (h reader) {
	return MakeStringReader(readerNew)
}

func readerNew(u lib.Umwelt, _ int, _ string) (z *lib.Zettel, err error) {
	id, err := u.Kasten.NewId()

	if err != nil {
		return
	}

	z = &lib.Zettel{
		Id:     id.Int(),
		Path:   lib.MakePathFromId(u.Kasten.BasePath(), id.String()),
		Umwelt: u,
		Note: lib.Note{
			Metadata: metadata.MakeMetadata(),
		},
	}

	z.Metadata.AddStringTags(u.TagsForNewZettels...)

	return
}
