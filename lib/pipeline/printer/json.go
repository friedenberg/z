package printer

import (
	"encoding/json"

	"github.com/friedenberg/z/lib"
	"github.com/friedenberg/z/util/stdprinter"
)

type JsonZettelPrinter struct{}

func (p *JsonZettelPrinter) Begin() {}
func (p *JsonZettelPrinter) End()   {}

func (p *JsonZettelPrinter) PrintZettel(i int, z *lib.Zettel, errIn error) {
	if errIn != nil {
		stdprinter.Err(errIn)
		return
	}

	b, errOut := json.Marshal(z.Metadata)

	if errOut != nil {
		stdprinter.Err(errOut)
		return
	}

	stdprinter.Out(string(b))
}
