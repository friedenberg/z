package printer

import (
	"encoding/json"

	"github.com/friedenberg/z/lib"
	"github.com/friedenberg/z/util"
)

type JsonZettelPrinter struct{}

func (p *JsonZettelPrinter) Begin() {}
func (p *JsonZettelPrinter) End()   {}

func (p *JsonZettelPrinter) PrintZettel(i int, z *lib.Zettel, errIn error) {
	if errIn != nil {
		util.StdPrinterErr(errIn)
		return
	}

	b, errOut := json.Marshal(z.IndexData)

	if errOut != nil {
		util.StdPrinterErr(errOut)
		return
	}

	util.StdPrinterOut(string(b))
}
