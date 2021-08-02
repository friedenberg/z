package printer

import (
	"strings"
	"sync"

	"github.com/friedenberg/z/lib"
	"github.com/friedenberg/z/util"
)

type AlfredJsonZettelPrinter struct {
	PrintFunc       func(i int, z *lib.Zettel, sb *strings.Builder)
	afterFirstPrint bool
	sync.Mutex
}

func (p *AlfredJsonZettelPrinter) Begin() {
	if p.PrintFunc == nil {
		p.PrintFunc = func(i int, z *lib.Zettel, sb *strings.Builder) {
			item := alfredItemFromZettelDefault(z)
			//TODO handle error
			j, _ := lib.GenerateAlfredJson(item)

			sb.WriteString(j)
		}
	}

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

func (p *AlfredJsonZettelPrinter) PrintZettel(i int, z *lib.Zettel, errIn error) {
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

	p.PrintFunc(i, z, &sb)
	util.StdPrinterOut(sb.String())
}

func (p *AlfredJsonZettelPrinter) End() {
	util.StdPrinterOut(`]}`)
}
