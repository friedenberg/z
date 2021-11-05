package commands

import (
	"github.com/friedenberg/z/lib/pipeline"
	"github.com/friedenberg/z/lib/zettel/reader"
	"golang.org/x/xerrors"
)

type attachmentKind struct {
	pipeline.Reader
}

func (a attachmentKind) String() string {
	//TODO-P4
	return ""
}

func (a *attachmentKind) Set(s string) (err error) {
	switch s {
	case "zettels-copy":
		a.Reader = reader.Import(true)
	case "zettels":
		a.Reader = reader.Import(false)
	case "files":
		a.Reader = reader.ForAttachmentFile()
	case "urls":
		a.Reader = reader.ForAttachmentUrl()
	default:
		err = xerrors.Errorf("unsupported type: '%s'", s)
		return
	}

	return
}
