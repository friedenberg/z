package reader

import (
	"github.com/friedenberg/z/lib"
	"github.com/friedenberg/z/lib/zettel"
)

func Make(f ReaderFunc) (h reader) {
	h.readerFunc = f
	return
}

func MakeStringReader(f StringReaderFunc) (h reader) {
	h.readerFunc = func(u lib.Umwelt, i int, b []byte) (*zettel.Zettel, error) {
		return f(u, i, string(b))
	}
	return
}
