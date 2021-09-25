package pipeline

import "github.com/friedenberg/z/lib"

type Hydrator func(int, string) (*lib.KastenZettel, error)

type Filter func(int, *lib.KastenZettel) bool

type Printer interface {
	Begin()
	PrintZettel(int, *lib.KastenZettel, error)
	End()
}

type FilterPrinter struct {
	Filter  Filter
	Printer Printer
}
