package pipeline

import (
	"io"

	"github.com/friedenberg/z/lib"
	"github.com/friedenberg/z/lib/zettel"
)

type Reader interface {
	ReadZettel(*lib.Umwelt, int, []byte) (*zettel.Zettel, error)
}

type Filter interface {
	FilterZettel(int, *zettel.Zettel) bool
}

type Modifier interface {
	ModifyZettel(int, *zettel.Zettel) error
}

type Writer interface {
	WriteZettel(io.Writer, int, *zettel.Zettel)
}

type WriterExcludeEmpty interface {
	Writer
	SetExcludeEmpty()
}

type WriterError interface {
	WriteZettelError(io.Writer, int, *zettel.Zettel, error)
}

