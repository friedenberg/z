package lib

import (
	"github.com/friedenberg/z/lib/kasten"
	"github.com/friedenberg/z/lib/zettel/metadata"
)

type Kasten struct {
	Local   *FilesAndGit
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
