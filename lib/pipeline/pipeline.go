package pipeline

import (
	"io"
	"os"
	"sync"

	"github.com/friedenberg/z/lib"
	"github.com/friedenberg/z/lib/feeder"
	"github.com/friedenberg/z/lib/zettel"
	"github.com/friedenberg/z/lib/zettel/beginner"
	"github.com/friedenberg/z/lib/zettel/ender"
	"github.com/friedenberg/z/lib/zettel/reader"
	"github.com/friedenberg/z/util/stdprinter"
	"golang.org/x/xerrors"
)

//TODO-P3 add structured error printing
type Pipeline struct {
	//TODO-P3 modify to support non-strings
	feeder.Feeder
	Reader
	Filter
	Modifier
	Writer
	Out io.Writer
}

func (p Pipeline) Run(u *lib.Umwelt) {
	wg := &sync.WaitGroup{}

	p.begin()
	defer p.end()

	p.runAll(u, wg)
}

func (p Pipeline) runAll(u *lib.Umwelt, wg *sync.WaitGroup) {
	defer wg.Wait()

	go p.Feeder.Run()

	i := 0

	for arg := range p.Feeder.GetChan() {
		wg.Add(1)
		go p.runOne(u, wg, i, arg)
		i += 1
	}
}

func (p Pipeline) runOne(u *lib.Umwelt, wg *sync.WaitGroup, i int, s string) {
	var err error
	defer func() {
		if err != nil {
			stdprinter.Error(err)
		}
	}()

	defer wg.Done()

	z, err := p.readZettel(u, i, s)

	if err != nil {
		err = xerrors.Errorf("failed to read zettel (%s, '%s': %w", i, s, err)
		return
	} else if z == nil {
		err = xerrors.Errorf("zettel reader returned nil for '%s'", s)
		return
	}

	if p.Filter != nil && !p.FilterZettel(i, z) {
		return
	}

	err = p.modifyZettel(i, z)

	if err != nil {
		err = xerrors.Errorf("failed to modify zettel '%s': %w", z.Id, err)
		return
	}

	p.writeZettel(i, z)
}

func (p Pipeline) begin() {
	if c, ok := p.Modifier.(beginner.Beginner); ok {
		c.Begin(p.outWriter())
	}

	if c, ok := p.Writer.(beginner.Beginner); ok {
		c.Begin(p.outWriter())
	}
}

func (p Pipeline) end() {
	if c, ok := p.Modifier.(ender.Ender); ok {
		c.End(p.outWriter())
	}

	if c, ok := p.Writer.(ender.Ender); ok {
		c.End(p.outWriter())
	}
}

func (p Pipeline) outWriter() (w io.Writer) {
	w = p.Out

	if w == nil {
		w = os.Stdout
	}

	return
}

func (p Pipeline) readZettel(u *lib.Umwelt, i int, s string) (z *zettel.Zettel, err error) {
	if p.Reader != nil {
		return p.ReadZettel(u, i, []byte(s))
	}

	z, err = reader.FromIndex(u, i, s)

	return
}

func (p Pipeline) modifyZettel(i int, z *zettel.Zettel) (err error) {
	if p.Modifier == nil {
		return
	}

	err = p.ModifyZettel(i, z)

	return
}

func (p Pipeline) writeZettel(i int, z *zettel.Zettel) {
	if p.Writer == nil {
		return
	}

	// w := makeSyncedWriter(p.outWriter())
	// defer w.Flush()
	p.WriteZettel(p.outWriter(), i, z)
}
