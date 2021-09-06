package commands

import (
	"fmt"
	"sync"

	"github.com/friedenberg/z/commands/printer"
	"github.com/friedenberg/z/lib"
	"github.com/friedenberg/z/util"
)

type ArgNormalizeFunc func(int, string) (string, error)
type HydrateFunc func(int, *lib.Zettel, string) error
type ActionFunc func(int, *lib.Zettel) (bool, error)

type Processor struct {
	kasten        *lib.FilesAndGit
	files         []string
	waitGroup     sync.WaitGroup
	argNormalizer ArgNormalizeFunc
	hydrator      HydrateFunc
	actioner      ActionFunc
	printer       printer.ZettelPrinter
}

func MakeProcessor(e *lib.FilesAndGit, files []string, zp printer.ZettelPrinter) (processor *Processor) {
	if len(files) == 0 {
		var err error
		files, err = e.GetAll()

		if err != nil {
			panic(err)
		}
	}

	processor = &Processor{
		kasten:  e,
		files:   files,
		printer: &printer.MultiplexingZettelPrinter{Printer: zp},
	}

	return
}

func (p *Processor) init() {
	if p.argNormalizer == nil {
		p.argNormalizer = func(_ int, path string) (normalizedArg string, err error) {
			normalizedArg, err = p.kasten.GetNormalizedPath(path)
			return
		}
	}

	if p.hydrator == nil {
		p.hydrator = func(_ int, z *lib.Zettel, path string) error {
			z.Path = path
			return z.Hydrate(false)
		}
	}
}

func (p *Processor) SetPrinter(printer printer.ZettelPrinter) {
	p.printer = printer
}

func (p *Processor) Run() (err error) {
	p.init()

	runRead := func() {
		for i, file := range p.files {
			go func(i int, f string) {
				defer p.waitGroup.Done()
				z, err := p.HydrateFile(i, f)

				if err != nil {
					err = fmt.Errorf("%s: failed to hydrate: %w", f, err)
					p.printer.PrintZettel(i, z, err)
					return
				}

				err = p.ActionZettel(i, z)

				if err != nil {
					err = fmt.Errorf("%s:\n\t%w", f, err)
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
	//TODO move to read site
	util.OpenFilesGuardInstance.Lock()
	defer util.OpenFilesGuardInstance.Unlock()

	z = &lib.Zettel{
		FilesAndGit: p.kasten,
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
