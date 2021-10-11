package printer

import (
	"strings"

	"github.com/friedenberg/z/lib"
	"github.com/friedenberg/z/util"
)

type FullZettelPrinter struct{}

func (p *FullZettelPrinter) Begin() {}
func (p *FullZettelPrinter) End()   {}

func (p *FullZettelPrinter) PrintZettel(_ int, z *lib.Zettel, errIn error) {
	if errIn != nil {
		util.StdPrinterErr(errIn)
		return
	}

	sb := &strings.Builder{}

	y, err := z.Note.Metadata.ToYAMLWithBoundary()

	if err != nil {
		util.StdPrinterError(err)
		return
	}

	sb.WriteString(y)
	sb.WriteString(z.Body)
	util.StdPrinterOutf(sb.String())
}
