package zettel

import (
	"github.com/friedenberg/z/lib/zettel/metadata"
)

type ZUmwelt interface {
	Dir() string
}

type Zettel struct {
	ZUmwelt
	Id
	Note
	Path string
}

type Note struct {
	metadata.Metadata
	Body string
}
