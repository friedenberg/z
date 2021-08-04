package lib

type Zettel struct {
	Kasten       *Kasten
	Path      string
	Id        int64
	IndexData ZettelIndexData
	Data      ZettelData
}

type ZettelData struct {
	MetadataYaml string
	Body         string
}
