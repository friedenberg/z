package writer

import (
	"io"
	"sync"

	"github.com/friedenberg/z/lib"
	"github.com/friedenberg/z/util/stdprinter"
)

type AlfredJson struct {
	ItemFunc        func(z *lib.Zettel) []lib.AlfredItem
	afterFirstPrint bool
	sync.Mutex
}

func (p *AlfredJson) Begin(w io.Writer) {
	if p.ItemFunc == nil {
		p.ItemFunc = alfredItemsFromZettelDefault
	}

	io.WriteString(w, `{"items":[`)
}

func (p *AlfredJson) shouldPrintComma() bool {
	return p.afterFirstPrint
}

func (p *AlfredJson) setShouldPrintComma() {
	p.afterFirstPrint = true
}

func (p *AlfredJson) WriteZettel(w io.Writer, i int, z *lib.Zettel) {

	items := p.ItemFunc(z)
	//TODO-P2 handle erro
	j, _ := lib.GenerateAlfredItemsJson(items)

	p.Lock()
	defer p.Unlock()
	defer p.setShouldPrintComma()

	if p.shouldPrintComma() {
		_, err := io.WriteString(w, ",")
		stdprinter.PanicIfError(err)
	}

	_, err := io.WriteString(w, j)
	stdprinter.PanicIfError(err)
}

func (p *AlfredJson) End(w io.Writer) {
	io.WriteString(w, `]}`)
}
