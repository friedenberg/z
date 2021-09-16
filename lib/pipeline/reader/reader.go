package reader

import (
	"github.com/friedenberg/z/lib"
	"github.com/friedenberg/z/util/stdprinter"
	"golang.org/x/xerrors"
)

type reader struct {
	readerFunc ReaderFunc
}

func (h reader) ReadZettel(u lib.Umwelt, i int, b []byte) (*lib.Zettel, error) {
	if h.readerFunc == nil {
		stdprinter.PanicIfError(xerrors.Errorf("no hydrator set"))
	}

	return h.readerFunc(u, i, b)
}
