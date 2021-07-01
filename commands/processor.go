package commands

import (
	"fmt"
	"os"
	"sync"

	"github.com/friedenberg/z/lib"
)

type ArgNormalizeFunc func(int, string) (string, error)
type HydrateFunc func(int, *lib.Zettel, string) error
type ActionFunc func(int, *lib.Zettel) error

type Processor struct {
	env                  Env
	files                []string
	waitGroup            sync.WaitGroup
	openFileGuardChannel chan struct{}
	writeWaitGroup       sync.WaitGroup
	argNormalizer        ArgNormalizeFunc
	hydrator             HydrateFunc
	actioner             ActionFunc
	putter               Putter
}

func MakeProcessor(e Env, files []string, putter Putter) (processor *Processor) {
	processor = &Processor{
		env:                  e,
		files:                files,
		openFileGuardChannel: make(chan struct{}, 240),
		putter:               putter,
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
			return z.HydrateFromFilePath()
		}
	}
}

func (p *Processor) Run() (err error) {
	p.init()

	runRead := func() {
		for i, file := range p.files {
			p.waitGroup.Add(1)
			go func(i int, f string) {
				defer p.waitGroup.Done()
				z, err := p.HydrateFile(i, f)

				if err != nil {
					err = fmt.Errorf("%s: failed to hydrate: %w", f, err)
					fmt.Fprintln(os.Stderr, err)
					return
				}

				defer lib.ZettelPoolInstance.Put(z)

				err = p.ActionZettel(i, z)

				if err != nil {
					err = fmt.Errorf("%s: %w", f, err)
					fmt.Fprintln(os.Stderr, err)
				}
			}(i, file)
		}

		p.waitGroup.Wait()
		close(p.putter.GetChannel())
	}

	go runRead()
	p.putter.Print()

	return nil
}

func (p *Processor) HydrateFile(i int, path string) (z *lib.Zettel, err error) {
	p.openFileGuardChannel <- struct{}{}
	defer func() { <-p.openFileGuardChannel }()

	z = lib.ZettelPoolInstance.Get()

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

	p.putter.GetChannel() <- z

	return
}
