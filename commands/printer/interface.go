package printer

import "github.com/friedenberg/z/lib"

type ZettelPrinter interface {
	Begin()
	PrintZettel(int, *lib.Zettel, error)
	End()
}
