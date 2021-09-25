package lib

import "github.com/friedenberg/z/lib/kasten"

type Kasten struct {
	Local   *FilesAndGit
	Remotes map[string]kasten.RemoteImplementation
}

type Zettel struct {
	*Umwelt

	Id int64
	Note

	Path string
	Data ZettelData
}

type Note struct {
	Metadata
	Body string
}

type ZettelData struct {
	MetadataYaml string
}
