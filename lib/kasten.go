package lib

import (
	"github.com/friedenberg/z/lib/feeder"
	"github.com/friedenberg/z/lib/zettel"
)

type Kasten interface {
	Init(*Umwelt, map[string]interface{}) (err error)
	BasePath() string
	GetAll() feeder.Feeder
	GetNormalizedPath(a string) (b string, err error)
	NewId() (id zettel.Id, err error)
	Hydrate(u *Umwelt, z *zettel.Zettel, includeBody bool) (err error)
	CommitTransaction(*Umwelt) error
}
