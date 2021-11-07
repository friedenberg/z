package reader

import (
	"github.com/friedenberg/z/lib"
	"github.com/friedenberg/z/lib/zettel"
)

type ReaderFunc func(*lib.Umwelt, int, []byte) (*zettel.Zettel, error)
