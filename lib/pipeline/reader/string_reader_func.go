package reader

import "github.com/friedenberg/z/lib"

type StringReaderFunc func(lib.Umwelt, int, string) (*lib.Zettel, error)
