package printer

import (
	"strings"
	"sync"

	"github.com/friedenberg/z/lib"
	"github.com/friedenberg/z/util"
)

type AlfredJsonZettelPrinter struct {
	ItemFunc        func(z *lib.Zettel) []lib.AlfredItem
	afterFirstPrint bool
	sync.Mutex
}

func (p *AlfredJsonZettelPrinter) Begin() {
	if p.ItemFunc == nil {
		p.ItemFunc = alfredItemsFromZettelDefault
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

	items := p.ItemFunc(z)
	//TODO-P2 handle erro
	j, _ := lib.GenerateAlfredItemsJson(items)

	sb.WriteString(j)
	util.StdPrinterOut(sb.String())
}

func (p *AlfredJsonZettelPrinter) End() {
	util.StdPrinterOut(`]}`)
}
