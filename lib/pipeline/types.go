package pipeline

import (
	"io"

	"github.com/friedenberg/z/lib"
)

type Reader interface {
	ReadZettel(lib.Umwelt, int, []byte) (*lib.Zettel, error)
}

type Filter interface {
	FilterZettel(int, *lib.Zettel) bool
}

type Modifier interface {
	ModifyZettel(int, *lib.Zettel) error
}

type Beginner interface {
	Begin(io.Writer)
}

type Writer interface {
	WriteZettel(io.Writer, int, *lib.Zettel)
}

type WriterError interface {
	WriteZettelError(io.Writer, int, *lib.Zettel, error)
}

type Ender interface {
	End(io.Writer)
}
