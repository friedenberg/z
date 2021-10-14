package printer

import (
	"sync"

	"github.com/friedenberg/z/commands/options"
	"github.com/friedenberg/z/lib"
	"github.com/friedenberg/z/util"
	"github.com/friedenberg/z/util/stdprinter"
	"golang.org/x/xerrors"
)

type ActionZettelPrinter struct {
	Umwelt      lib.Umwelt
	Actions     options.Actions
	zettels     []*lib.Zettel
	zettelFiles []string
	files       []string
	urls        []string
}

func (p *ActionZettelPrinter) Begin() {
	p.zettels = make([]*lib.Zettel, 0)
	p.zettelFiles = make([]string, 0)
	p.files = make([]string, 0)
}

func (p *ActionZettelPrinter) PrintZettel(i int, z *lib.Zettel, errIn error) {
	if errIn != nil {
		stdprinter.Err(errIn)
		return
	}

	p.zettels = append(p.zettels, z)
	p.zettelFiles = append(p.zettelFiles, z.Path)

	if f, ok := z.Note.Metadata.LocalFile(); ok {
		p.files = append(p.files, f.FilePath(p.Umwelt.BasePath))
	}

	if u, ok := z.Note.Metadata.Url(); ok {
		p.urls = append(p.urls, u.String())
	}

	if p.Actions&options.ActionPrintZettelPath != 0 {
		//TODO-P2 full path
		stdprinter.Out(z.Path)
	}
}

func (p *ActionZettelPrinter) End() {
	wg := &sync.WaitGroup{}

	var err error

	runAction := func(actionFunc func() error) {
		wg.Add(1)

		go func() {
			defer wg.Done()

			err := actionFunc()

			if err != nil {
				stdprinter.Err(err)
				return
			}
		}()
	}

	if p.Actions&options.ActionEdit != 0 {
		runAction(p.openZettels)
	}

	if p.Actions&options.ActionOpenFile != 0 {
		runAction(p.openFiles)
	}

	if p.Actions&options.ActionOpenUrl != 0 {
		runAction(p.openUrls)
	}

	wg.Wait()

	for i, z := range p.zettels {
		p.Umwelt.Mod.PrintZettel(i, z, nil)
	}

	if err != nil {
		stdprinter.Err(err)
	}
}

func (p *ActionZettelPrinter) openZettels() (err error) {
	if len(p.zettels) == 0 {
		stdprinter.Debug("no zettels to open")
		return
	}

	args := []string{"-f", "-p"}

	cmd := util.ExecCommand(
		"mvim",
		args,
		p.zettelFiles,
	)

	output, err := cmd.CombinedOutput()

	if err != nil {
		err = xerrors.Errorf("opening zettels ('%q'): %s", p.zettels, output)
		return
	}

	return
}

func (p *ActionZettelPrinter) openFiles() (err error) {
	if len(p.files) == 0 {
		return
	}

	cmd := util.ExecCommand(
		"open",
		[]string{"-W"},
		p.files,
	)

	output, err := cmd.CombinedOutput()

	if err != nil {
		err = xerrors.Errorf("opening files ('%q'): %s", p.files, output)
		return
	}

	return
}

func (p *ActionZettelPrinter) openUrls() (err error) {
	if len(p.urls) == 0 {
		return
	}

	args := []string{
		"-na",
		"Google Chrome",
		"--args",
		"--new-window",
	}

	cmd := util.ExecCommand(
		"open",
		args,
		p.urls,
	)

	output, err := cmd.CombinedOutput()

	if err != nil {
		err = xerrors.Errorf("opening urls ('%q'): %s", p.urls, output)
		return
	}

	return
}
