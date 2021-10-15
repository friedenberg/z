package pipeline

import (
	"github.com/friedenberg/z/lib"
	"github.com/friedenberg/z/lib/zettel/metadata"
)

func New(u lib.Umwelt) (z *lib.Zettel, err error) {
	id, err := u.Store().NewId()

	if err != nil {
		return
	}

	z = &lib.Zettel{
		Id:     id.Int(),
		Path:   lib.MakePathFromId(u.Store().BasePath(), id.String()),
		Umwelt: u,
		Note: lib.Note{
			Metadata: metadata.MakeMetadata(),
		},
	}

	z.Metadata.AddStringTags(u.TagsForNewZettels...)

	return
}
