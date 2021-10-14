package printer

import (
	"github.com/friedenberg/z/lib"
	"github.com/friedenberg/z/util/stdprinter"
)

type FilenameZettelPrinter struct{}

func (p *FilenameZettelPrinter) Begin() {}
func (p *FilenameZettelPrinter) End()   {}

func (p *FilenameZettelPrinter) PrintZettel(i int, z *lib.Zettel, errIn error) {
	if errIn != nil {
		stdprinter.Err(errIn)
		return
	}

	stdprinter.Out(z.Path)
}
