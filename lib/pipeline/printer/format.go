package printer

import (
	"github.com/friedenberg/z/lib"
	"github.com/friedenberg/z/util/stdprinter"
)

type FormatZettelPrinter struct {
	Formatter lib.Formatter
}

func (p *FormatZettelPrinter) Begin() {}
func (p *FormatZettelPrinter) End()   {}

func (p *FormatZettelPrinter) PrintZettel(i int, z *lib.Zettel, errIn error) {
	if errIn != nil {
		stdprinter.Err(errIn)
		return
	}

	//TODO-P4 should empty strings be printed?
	stdprinter.Outf("%s", p.Formatter.Format(z))
}
