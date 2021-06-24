package commands

import (
	"path/filepath"
	"sync"

	"github.com/friedenberg/z/lib"
)

type ProcessorAction func(z *lib.Zettel)

type Processor struct {
	files                []string
	waitGroup            sync.WaitGroup
	openFileGuardChannel chan struct{}
	writeWaitGroup       sync.WaitGroup
	parallelAction       ProcessorAction
	putter               Putter
}

func MakeProcessor(glob string, pa ProcessorAction, putter Putter) (processor *Processor, err error) {
	files, err := filepath.Glob(glob)

	if err != nil {
		return
	}

	processor = &Processor{
		files:                files,
		openFileGuardChannel: make(chan struct{}, 240),
		putter:               putter,
		parallelAction:       pa,
	}

	return
}

func (p *Processor) Run() (err error) {
	runRead := func() {
		for _, file := range p.files {
			p.waitGroup.Add(1)
			go p.ProcessFile(file)
		}

		p.waitGroup.Wait()
		close(p.putter.GetChannel())
	}

	go runRead()
	p.putter.Print()

	return nil
}

func (p *Processor) ProcessFile(path string) {
	defer p.waitGroup.Done()

	p.openFileGuardChannel <- struct{}{}
	defer func() { <-p.openFileGuardChannel }()

	z := lib.ZettelPoolInstance.Get()
	defer lib.ZettelPoolInstance.Put(z)

	z.HydrateFromFilePath(path)
	if p.parallelAction != nil {
		p.parallelAction(z)
	}

	p.putter.GetChannel() <- z
}
