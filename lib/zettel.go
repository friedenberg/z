package lib

type Zettel struct {
	Env        *Env
	Path       string
	IndexData  ZettelIndexData
	Data       ZettelData
	AlfredData AlfredData
}

type ZettelData struct {
	MetadataYaml string
	Body         string
}

type AlfredData struct {
	Item ZettelAlfredItem
	Json string
}
