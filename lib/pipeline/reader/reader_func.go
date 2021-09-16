package reader

import "github.com/friedenberg/z/lib"

type ReaderFunc func(lib.Umwelt, int, []byte) (*lib.Zettel, error)
