package printer

import (
	"strings"

	"github.com/friedenberg/z/lib"
	"github.com/friedenberg/z/util"
)

type FullZettelPrinter struct{}

func (p *FullZettelPrinter) Begin() {}
func (p *FullZettelPrinter) End()   {}

func (p *FullZettelPrinter) PrintZettel(_ int, z *lib.KastenZettel, errIn error) {
	if errIn != nil {
		util.StdPrinterErr(errIn)
		return
	}

	sb := &strings.Builder{}
	sb.WriteString(lib.MetadataStartSequence)
	sb.WriteString(z.Data.MetadataYaml)
	sb.WriteString(lib.MetadataEndSequence)
	sb.WriteString(z.Body)
	util.StdPrinterOutf(sb.String())
}
