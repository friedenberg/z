package zettel

import (
	"github.com/friedenberg/z/lib/zettel/metadata"
)

type ZUmwelt interface {
	Dir() string
}

type Zettel struct {
	ZUmwelt

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
