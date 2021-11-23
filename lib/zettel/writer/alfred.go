package writer

import (
	"io"
	"sync"

	"github.com/friedenberg/z/lib"
	"github.com/friedenberg/z/lib/zettel"
	"github.com/friedenberg/z/util/stdprinter"
	"golang.org/x/xerrors"
)

type AlfredJson struct {
	ItemFunc        func(z *zettel.Zettel) []lib.AlfredItem
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

func (p *AlfredJson) WriteZettel(w io.Writer, i int, z *zettel.Zettel) {
	items := p.ItemFunc(z)
	j, err := lib.GenerateAlfredItemsJson(items)

	if err != nil {
		err = xerrors.Errorf("failed to generated alfred items: %w", err)
		stdprinter.Error(err)
	}

	p.Lock()
	defer p.Unlock()
	defer p.setShouldPrintComma()

	if p.shouldPrintComma() {
		_, err := io.WriteString(w, ",")
		stdprinter.PanicIfError(err)
	}

	_, err = io.WriteString(w, j)
	stdprinter.PanicIfError(err)
}

func (p *AlfredJson) End(w io.Writer) {
	io.WriteString(w, `]}`)
}
