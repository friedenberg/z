package printer

import (
	"github.com/friedenberg/z/lib"
	"github.com/friedenberg/z/util"
)

type NullZettelPrinter struct{}

func (p *NullZettelPrinter) Begin() {}
func (p *NullZettelPrinter) End()   {}

func (p *NullZettelPrinter) PrintZettel(_ int, _ *lib.Zettel, errIn error) {
	if errIn != nil {
		util.StdPrinterErr(errIn)
		return
	}
}
