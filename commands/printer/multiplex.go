package printer

import (
	"sync"

	"github.com/friedenberg/z/lib"
	"github.com/friedenberg/z/lib/pipeline"
)

type multiplexPrintLine struct {
	i int
	z *lib.Zettel
	e error
}

type MultiplexingZettelPrinter struct {
	Printer   pipeline.Printer
	channel   chan multiplexPrintLine
	waitGroup *sync.WaitGroup
}

func (p *MultiplexingZettelPrinter) Begin() {
	p.channel = make(chan multiplexPrintLine)
	p.Printer.Begin()
	p.waitGroup = &sync.WaitGroup{}
	p.waitGroup.Add(1)

	go func() {
		defer p.waitGroup.Done()

		for l := range p.channel {
			p.Printer.PrintZettel(l.i, l.z, l.e)
		}
	}()
}

func (p *MultiplexingZettelPrinter) PrintZettel(i int, z *lib.Zettel, e error) {
	p.channel <- multiplexPrintLine{i, z, e}
}

func (p *MultiplexingZettelPrinter) End() {
	close(p.channel)
	p.waitGroup.Wait()
	p.Printer.End()
}
