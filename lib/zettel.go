package lib

import (
	"github.com/friedenberg/z/lib/zettel"
	"github.com/friedenberg/z/lib/zettel/metadata"
)

//TODO-P3 rename to LocalStore
type Kasten interface {
	Init(Umwelt, map[string]interface{}) (err error)
	BasePath() string
	GetAll() (zettels []string, err error)
	GetNormalizedPath(a string) (b string, err error)
	NewId() (id zettel.Id, err error)
	Hydrate(z *Zettel, includeBody bool) (err error)
	CommitTransaction(Umwelt) error
}

type Zettel struct {
	Umwelt

	//TODO-P2 change to zettel.Id
	Id int64
	Note

	Path string
}

type Note struct {
	// Metadata
	Metadata metadata.Metadata
	Body     string
}
