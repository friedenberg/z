package lib

type Zettel struct {
	Kasten   *Kasten
	Id       int64
	Metadata Metadata
	Body     string

	Path string
	Data ZettelData
}

type ZettelData struct {
	MetadataYaml string
}
