package lib

import (
	"fmt"
	"strings"

	"github.com/friedenberg/z/lib/zettel/metadata"
)

func (z *Zettel) Merge(z1 *Zettel) (err error) {
	z.Note = z.Note.Merge(z1.Note)
	return
}

func (n Note) Merge(n1 Note) (n2 Note) {
	m := metadata.MakeMetadata()
	m.AddStringTags(n1.Metadata.TagStrings()...)
	//intentionally clobber n1 with n's tags to preserve them
	m.AddStringTags(n.Metadata.TagStrings()...)

	m.SetDescription(
		fmt.Sprintf(
			"%s %s",
			n.Metadata.Description(),
			n1.Metadata.Description(),
		),
	)

	n2.Metadata = m
	n2.Body = strings.TrimSpace(fmt.Sprintf("%s\n\n%s", n.Body, n1.Body))
	return
}
