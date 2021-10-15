package commands

import (
	"sync"

	"github.com/friedenberg/z/lib"
	"github.com/friedenberg/z/lib/pipeline"
	"github.com/friedenberg/z/lib/pipeline/printer"
	"github.com/friedenberg/z/util"
	"golang.org/x/xerrors"
)

type ArgNormalizeFunc func(int, string) (string, error)
type HydrateFunc func(int, *lib.Zettel, string) error
type ActionFunc func(int, *lib.Zettel) (bool, error)

type Processor struct {
	kasten        lib.Umwelt
	files         []string
	waitGroup     sync.WaitGroup
	argNormalizer ArgNormalizeFunc
	hydrator      HydrateFunc
	actioner      ActionFunc
	printer       pipeline.Printer
}

func MakeProcessor(e lib.Umwelt, files []string, zp pipeline.Printer) (processor *Processor) {
	processor = &Processor{
		kasten:  e,
		files:   files,
		printer: &printer.MultiplexingZettelPrinter{Printer: zp},
	}

	return
}

func (p *Processor) init() (err error) {
	if len(p.files) == 0 {
		p.files, err = p.getDefaultFiles()

		if err != nil {
			return
		}
	}

	if p.argNormalizer == nil {
		p.argNormalizer = DefaultArgNormalizer(p.kasten)
	}

	if p.hydrator == nil {
		p.hydrator = HydrateFromIndexFunc(p.kasten)
	}

	return
}

func (p *Processor) getDefaultFiles() (files []string, err error) {
	files = p.kasten.GetAll()

	return
}

func (p *Processor) SetPrinter(printer pipeline.Printer) {
	p.printer = printer
}

func (p *Processor) Run() (err error) {
	err = p.init()

	if err != nil {
		return
	}

	iter := func(i int, f string) (iterErr error) {
		z, iterErr := p.HydrateFile(i, f)

		if iterErr != nil {
			iterErr = xerrors.Errorf("%s: failed to hydrate: %w", f, iterErr)
			return
		}

		iterErr = p.ActionZettel(i, z)

		if iterErr != nil {
			return
			iterErr = xerrors.Errorf("%s:\n\t%w", f, iterErr)
		}

		p.printer.PrintZettel(i, z, nil)

		return
	}

	errIter := func(i int, s string, err error) {
		p.printer.PrintZettel(i, nil, err)
	}

	p.printer.Begin()
	par := util.Parallelizer{Args: p.files}
	par.Run(iter, errIter)
	p.printer.End()

	return nil
}

func (p *Processor) HydrateFile(i int, path string) (z *lib.Zettel, err error) {
	z = &lib.Zettel{
		Umwelt: p.kasten,
	}

	a, err := p.argNormalizer(i, path)

	if err != nil {
		return
	}

	err = p.hydrator(i, z, a)

	return
}

func (p *Processor) ActionZettel(i int, z *lib.Zettel) (err error) {
	shouldPrint := true

	if p.actioner != nil {
		shouldPrint, err = p.actioner(i, z)
	}

	if err != nil {
		return
	}

	if shouldPrint {
		p.printer.PrintZettel(i, z, nil)
	}

	return
}
