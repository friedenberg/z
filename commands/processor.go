package commands

import (
	"fmt"
	"sync"

	"github.com/friedenberg/z/lib"
	"github.com/friedenberg/z/util"
)

type ArgNormalizeFunc func(int, string) (string, error)
type HydrateFunc func(int, *lib.Zettel, string) error
type ActionFunc func(int, *lib.Zettel) error

type Processor struct {
	env           *lib.Env
	files         []string
	waitGroup     sync.WaitGroup
	argNormalizer ArgNormalizeFunc
	hydrator      HydrateFunc
	actioner      ActionFunc
	printer       zettelPrinter
}

func MakeProcessor(e *lib.Env, files []string, zp zettelPrinter) (processor *Processor) {
	if len(files) == 0 {
		var err error
		files, err = e.GetAllZettels()

		if err != nil {
			panic(err)
		}
	}

	processor = &Processor{
		env:     e,
		files:   files,
		printer: &multiplexingZettelPrinter{printer: zp},
	}

	return
}

func (p *Processor) init() {
	if p.argNormalizer == nil {
		p.argNormalizer = func(_ int, path string) (normalizedArg string, err error) {
			normalizedArg, err = p.env.GetNormalizedPath(path)
			return
		}
	}

	if p.hydrator == nil {
		p.hydrator = func(_ int, z *lib.Zettel, path string) error {
			z.Path = path
			return z.HydrateFromFilePath(false)
		}
	}
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
					p.printer.printZettel(i, z, err)
					return
				}

				defer p.env.ZettelPool.Put(z)

				err = p.ActionZettel(i, z)

				if err != nil {
					err = fmt.Errorf("%s:\n\t%w", f, err)
					p.printer.printZettel(i, z, err)
				}
			}(i, file)
		}
	}

	p.waitGroup.Add(len(p.files))

	p.printer.begin()

	go runRead()
	p.waitGroup.Wait()
	p.printer.end()

	return nil
}

func (p *Processor) HydrateFile(i int, path string) (z *lib.Zettel, err error) {
	//TODO move to read site
	util.OpenFilesGuardInstance.Lock()
	defer util.OpenFilesGuardInstance.Unlock()

	z = p.env.ZettelPool.Get()

	a, err := p.argNormalizer(i, path)

	if err != nil {
		return
	}

	err = p.hydrator(i, z, a)

	return
}

func (p *Processor) ActionZettel(i int, z *lib.Zettel) (err error) {
	if p.actioner != nil {
		err = p.actioner(i, z)
	}

	if err != nil {
		return
	}

	p.printer.printZettel(i, z, nil)

	return
}
