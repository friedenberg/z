package printer

import (
	"strings"
	"sync"

	"github.com/friedenberg/z/lib"
	"github.com/friedenberg/z/util"
)

type AlfredJsonZettelPrinter struct {
	afterFirstPrint bool
	sync.Mutex
}

func (p *AlfredJsonZettelPrinter) Begin() {
	util.StdPrinterOut(`{"items":[`)
}

func (p *AlfredJsonZettelPrinter) shouldPrintComma() bool {
	p.Lock()
	defer p.Unlock()

	return p.afterFirstPrint
}

func (p *AlfredJsonZettelPrinter) setShouldPrintComma() {
	p.Lock()
	defer p.Unlock()

	p.afterFirstPrint = true
}

func (p *AlfredJsonZettelPrinter) PrintZettel(_ int, z *lib.Zettel, errIn error) {
	defer p.setShouldPrintComma()

	if errIn != nil {
		util.StdPrinterErr(errIn)
		return
	}

	sb := strings.Builder{}
	if p.shouldPrintComma() {
		sb.WriteString(",")
		sb.WriteString("\n")
	}

	sb.WriteString(z.AlfredData.Json)
	util.StdPrinterOut(sb.String())
}

func (p *AlfredJsonZettelPrinter) End() {
	util.StdPrinterOut(`]}`)
}
