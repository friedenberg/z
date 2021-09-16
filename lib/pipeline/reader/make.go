package reader

import "github.com/friedenberg/z/lib"

func Make(f ReaderFunc) (h reader) {
	h.readerFunc = f
	return
}

func MakeStringReader(f StringReaderFunc) (h reader) {
	h.readerFunc = func(u lib.Umwelt, i int, b []byte) (*lib.Zettel, error) {
		return f(u, i, string(b))
	}
	return
}
