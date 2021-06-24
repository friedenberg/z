package lib

type Zettel struct {
	Path           string
	MetadataYaml   string
	Body           string
	Metadata       ZettelMetadata
	AlfredItem     ZettelAlfredItem
	AlfredItemJson string
}
