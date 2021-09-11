package commands

import (
	"sync"

	"github.com/friedenberg/z/commands/printer"
	"github.com/friedenberg/z/lib"
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
	printer       printer.ZettelPrinter
}

func MakeProcessor(e lib.Umwelt, files []string, zp printer.ZettelPrinter) (processor *Processor) {
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
		if p.kasten.Config.UseIndexCache {
			p.hydrator = HydrateFromIndexFunc(p.kasten)
		} else {
			p.hydrator = HydrateFromFileFunc(p.kasten, false)
		}
	}

	return
}

func (p *Processor) getDefaultFiles() (files []string, err error) {
	if p.kasten.Config.UseIndexCache {
		files = p.kasten.GetAll()
	} else {
		files, err = p.kasten.FilesAndGit().GetAll()
	}

	return
}

func (p *Processor) SetPrinter(printer printer.ZettelPrinter) {
	p.printer = printer
}

func (p *Processor) Run() (err error) {
	err = p.init()

	if err != nil {
		return
	}

	runRead := func() {
		for i, file := range p.files {
			go func(i int, f string) {
				defer p.waitGroup.Done()
				z, err := p.HydrateFile(i, f)

				if err != nil {
					err = xerrors.Errorf("%s: failed to hydrate: %w", f, err)
					p.printer.PrintZettel(i, z, err)
					return
				}

				err = p.ActionZettel(i, z)

				if err != nil {
					err = xerrors.Errorf("%s:\n\t%w", f, err)
					p.printer.PrintZettel(i, z, err)
				}
			}(i, file)
		}
	}

	p.waitGroup.Add(len(p.files))

	p.printer.Begin()

	go runRead()
	p.waitGroup.Wait()
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
