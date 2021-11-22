package reader

import (
	"github.com/friedenberg/z/lib"
	"github.com/friedenberg/z/lib/zettel"
	"github.com/friedenberg/z/lib/zettel/metadata"
)

func New() (h reader) {
	return MakeStringReader(readerNew)
}

func readerNew(u *lib.Umwelt, _ int, _ string) (z *zettel.Zettel, err error) {
	id, err := u.Kasten.NewId()

	if err != nil {
		return
	}

	z = &zettel.Zettel{
		Id:      id,
		Path:    lib.MakePathFromId(u.Kasten.BasePath(), id.String()),
		ZUmwelt: u,
		Note: zettel.Note{
			Metadata: metadata.MakeMetadata(),
		},
	}

	z.Metadata.AddStringTags(u.TagsForNewZettels...)

	return
}
