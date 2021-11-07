package reader

import (
	"github.com/friedenberg/z/lib"
	"github.com/friedenberg/z/lib/zettel"
)

type StringReaderFunc func(*lib.Umwelt, int, string) (*zettel.Zettel, error)
