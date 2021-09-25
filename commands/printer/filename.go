package printer

import (
	"github.com/friedenberg/z/lib"
	"github.com/friedenberg/z/util"
)

type FilenameZettelPrinter struct{}

func (p *FilenameZettelPrinter) Begin() {}
func (p *FilenameZettelPrinter) End()   {}

func (p *FilenameZettelPrinter) PrintZettel(i int, z *lib.KastenZettel, errIn error) {
	if errIn != nil {
		util.StdPrinterErr(errIn)
		return
	}

	util.StdPrinterOut(z.Path)
}
