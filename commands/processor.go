package commands

import (
	"fmt"
	"os"
	"sync"

	"github.com/friedenberg/z/lib"
)

type ProcessorAction func(i int, z *lib.Zettel) error

type Processor struct {
	env                  Env
	files                []string
	waitGroup            sync.WaitGroup
	openFileGuardChannel chan struct{}
	writeWaitGroup       sync.WaitGroup
	hydrateAction        ProcessorAction
	parallelAction       ProcessorAction
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

func (p *Processor) Run() (err error) {
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

	path, err = p.env.GetNormalizedPath(path)

	if err != nil {
		fmt.Println(err)
		return
	}

	z.Path = path

	if p.hydrateAction == nil {
		err = z.HydrateFromFilePath()
	} else {
		err = p.hydrateAction(i, z)
	}

	return
}

func (p *Processor) ActionZettel(i int, z *lib.Zettel) (err error) {
	if p.parallelAction != nil {
		err = p.parallelAction(i, z)
	}

	if err != nil {
		return
	}

	p.putter.GetChannel() <- z

	return
}
