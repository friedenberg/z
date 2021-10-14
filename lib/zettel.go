package lib

import (
	"github.com/friedenberg/z/lib/kasten"
	"github.com/friedenberg/z/lib/zettel"
	"github.com/friedenberg/z/lib/zettel/metadata"
)

type Store interface {
	BasePath() string
	GetAll() (zettels []string, err error)
	GetNormalizedPath(a string) (b string, err error)
	NewId() (id zettel.Id, err error)
	CommitTransaction(Umwelt) error
}

type Kasten struct {
	Local   Store
	Remotes map[string]kasten.RemoteImplementation
}

type Zettel struct {
	*Umwelt

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
