package lib

import "github.com/friedenberg/z/lib/kasten"

type Kasten struct {
	*Umwelt
	kasten.LocalImplementation
}

type KastenZettel struct {
	Kasten
	Zettel
}

type Zettel struct {
	Id       int64
	Metadata Metadata
	Body     string

	Path string
	Data ZettelData
}

type ZettelData struct {
	MetadataYaml string
}
