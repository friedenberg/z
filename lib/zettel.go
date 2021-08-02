package lib

type Zettel struct {
	Env       *Env
	Path      string
	Id        int64
	IndexData ZettelIndexData
	Data      ZettelData
}

type ZettelData struct {
	MetadataYaml string
	Body         string
}
