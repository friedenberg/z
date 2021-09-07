package lib

type Zettel struct {
	Umwelt   Umwelt
	Id       int64
	Metadata Metadata
	Body     string

	Path string
	Data ZettelData
}

type ZettelData struct {
	MetadataYaml string
}
