package commands

import (
	"encoding/json"
	"strings"
	"sync"

	"github.com/friedenberg/z/lib"
	"github.com/friedenberg/z/util"
)

type zettelPrinter interface {
	begin()
	printZettel(int, *lib.Zettel, error)
	end()
}

//   _   _       _ _
//  | \ | |_   _| | |
//  |  \| | | | | | |
//  | |\  | |_| | | |
//  |_| \_|\__,_|_|_|
//

type nullZettelPrinter struct{}

func (p *nullZettelPrinter) begin()                                    {}
func (p *nullZettelPrinter) printZettel(_ int, _ *lib.Zettel, _ error) {}
func (p *nullZettelPrinter) end()                                      {}

type multiplexPrintLine struct {
	i int
	z *lib.Zettel
	e error
}

type multiplexingZettelPrinter struct {
	printer zettelPrinter
	channel chan multiplexPrintLine
}

func (p *multiplexingZettelPrinter) begin() {
	p.channel = make(chan multiplexPrintLine)
	p.printer.begin()

	go func() {
		for l := range p.channel {
			p.printer.printZettel(l.i, l.z, l.e)
		}
	}()
}

func (p *multiplexingZettelPrinter) printZettel(i int, z *lib.Zettel, e error) {
	p.channel <- multiplexPrintLine{i, z, e}
}

func (p *multiplexingZettelPrinter) end() {
	close(p.channel)
	p.printer.end()
}

//   _____ _ _
//  |  ___(_) | ___ _ __   __ _ _ __ ___   ___
//  | |_  | | |/ _ \ '_ \ / _` | '_ ` _ \ / _ \
//  |  _| | | |  __/ | | | (_| | | | | | |  __/
//  |_|   |_|_|\___|_| |_|\__,_|_| |_| |_|\___|
//

type filenameZettelPrinter struct{}

func (p *filenameZettelPrinter) begin() {}
func (p *filenameZettelPrinter) end()   {}

func (p *filenameZettelPrinter) printZettel(i int, z *lib.Zettel, errIn error) {
	//TODO handle errIn
	util.StdPrinterOut(z.Path)
}

//       _ ____   ___  _   _
//      | / ___| / _ \| \ | |
//   _  | \___ \| | | |  \| |
//  | |_| |___) | |_| | |\  |
//   \___/|____/ \___/|_| \_|
//

type jsonZettelPrinter struct{}

func (p *jsonZettelPrinter) begin() {}
func (p *jsonZettelPrinter) end()   {}

func (p *jsonZettelPrinter) printZettel(i int, z *lib.Zettel, errIn error) {
	//TODO handle errIn
	b, errOut := json.Marshal(z.IndexData)

	if errOut != nil {
		util.StdPrinterErr(errOut)
		return
	}

	util.StdPrinterOut(string(b))
}

//      _    _  __              _
//     / \  | |/ _|_ __ ___  __| |
//    / _ \ | | |_| '__/ _ \/ _` |
//   / ___ \| |  _| | |  __/ (_| |
//  /_/   \_\_|_| |_|  \___|\__,_|
//

type alfredJsonZettelPrinter struct {
	afterFirstPrint bool
	sync.Mutex
}

func (p *alfredJsonZettelPrinter) begin() {
	util.StdPrinterOutf(`{"items":[`)
}

func (p *alfredJsonZettelPrinter) shouldPrintComma() bool {
	p.Lock()
	defer p.Unlock()

	return p.afterFirstPrint
}

func (p *alfredJsonZettelPrinter) setShouldPrintComma() {
	p.Lock()
	defer p.Unlock()

	p.afterFirstPrint = true
}

func (p *alfredJsonZettelPrinter) printZettel(_ int, z *lib.Zettel, _ error) {
	defer p.setShouldPrintComma()
	//TODO handle error
	sb := strings.Builder{}
	if p.shouldPrintComma() {
		sb.WriteString(",")
		sb.WriteString("\n")
	}

	sb.WriteString(z.AlfredData.Json)
	util.StdPrinterOutf(sb.String())
}

func (p *alfredJsonZettelPrinter) end() {
	util.StdPrinterOutf(`]}`)
}

//   _____      _ _
//  |  ___|   _| | |
//  | |_ | | | | | |
//  |  _|| |_| | | |
//  |_|   \__,_|_|_|
//

type fullZettelPrinter struct{}

func (p *fullZettelPrinter) begin() {}
func (p *fullZettelPrinter) end()   {}

func (p *fullZettelPrinter) printZettel(_ int, z *lib.Zettel, err error) {
	//todo handle error
	sb := &strings.Builder{}
	sb.WriteString(z.Data.MetadataYaml)
	sb.WriteString("\n")
	sb.WriteString(z.Data.Body)
	util.StdPrinterOutf(sb.String())
}
