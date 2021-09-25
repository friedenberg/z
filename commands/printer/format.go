package printer

import (
	"github.com/friedenberg/z/lib"
	"github.com/friedenberg/z/util"
)

type FormatZettelPrinter struct {
	Formatter lib.Formatter
}

func (p *FormatZettelPrinter) Begin() {}
func (p *FormatZettelPrinter) End()   {}

func (p *FormatZettelPrinter) PrintZettel(i int, z *lib.Zettel, errIn error) {
	if errIn != nil {
		util.StdPrinterErr(errIn)
		return
	}

	//TODO should empty strings be printed?
	util.StdPrinterOutf("%s", p.Formatter.Format(z))
}
