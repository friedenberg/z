package lib

import "github.com/friedenberg/z/lib/zettel"

type Kasten interface {
	Init(Umwelt, map[string]interface{}) (err error)
	BasePath() string
	GetAll() (zettels []string, err error)
	GetNormalizedPath(a string) (b string, err error)
	NewId() (id zettel.Id, err error)
	Hydrate(z *zettel.Zettel, includeBody bool) (err error)
	CommitTransaction(Umwelt) error
}
