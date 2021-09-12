package pipeline

import "github.com/friedenberg/z/lib"

type Filter func(int, *lib.Zettel) bool

type Printer interface {
	Begin()
	PrintZettel(int, *lib.Zettel, error)
	End()
}
