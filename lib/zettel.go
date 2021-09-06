package lib

type Zettel struct {
	FilesAndGit *FilesAndGit
	Id          int64
	Metadata    Metadata
	Body        string

	Path string
	Data ZettelData
}

type ZettelData struct {
	MetadataYaml string
}
