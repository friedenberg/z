package printer

import (
	"sync"

	"github.com/friedenberg/z/commands/options"
	"github.com/friedenberg/z/lib"
	"github.com/friedenberg/z/util"
	"golang.org/x/xerrors"
)

type ActionZettelPrinter struct {
	Umwelt      lib.Umwelt
	Actions     options.Actions
	zettels     []*lib.Zettel
	zettelFiles util.GitAnnex
	files       util.GitAnnex
	urls        []string
}

func (p *ActionZettelPrinter) Begin() {
	git := util.Git{
		Path: p.Umwelt.Kasten.Local.BasePath,
	}

	p.zettelFiles = util.GitAnnex{
		GitFilesToCommit: util.GitFilesToCommit{
			Git:   git,
			Files: make([]string, 0),
		},
	}

	p.files = util.GitAnnex{
		GitFilesToCommit: util.GitFilesToCommit{
			Git:   git,
			Files: make([]string, 0),
		},
	}
}

func (p *ActionZettelPrinter) PrintZettel(i int, z *lib.Zettel, errIn error) {
	if errIn != nil {
		util.StdPrinterErr(errIn)
		return
	}

	p.zettels = append(p.zettels, z)
	p.zettelFiles.Files = append(p.zettelFiles.Files, z.Path)

	if f, ok := z.Note.Metadata.LocalFile(); ok {
		p.files.Files = append(p.files.Files, f.FilePath(p.Umwelt.BasePath))
	}

	if u, ok := z.Note.Metadata.Url(); ok {
		p.urls = append(p.urls, u.String())
	}

	if p.Actions&options.ActionPrintZettelPath != 0 {
		//TODO-P2 full path
		util.StdPrinterOut(z.Path)
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
				util.StdPrinterErr(err)
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
		util.StdPrinterErr(err)
	}
}

func (p *ActionZettelPrinter) openZettels() (err error) {
	if len(p.zettels) == 0 {
		return
	}

	if p.Umwelt.Kasten.Local.GitAnnexEnabled {
		err = p.zettelFiles.Unlock()

		if err != nil {
			return
		}

		defer p.zettelFiles.Lock()
	}

	args := []string{"-f", "-p"}

	cmd := util.ExecCommand(
		"mvim",
		args,
		p.zettelFiles.Files,
	)

	output, err := cmd.CombinedOutput()

	if err != nil {
		err = xerrors.Errorf("opening zettels ('%q'): %s", p.zettels, output)
		return
	}

	return
}

func (p *ActionZettelPrinter) openFiles() (err error) {
	if len(p.files.Files) == 0 {
		return
	}

	if p.Umwelt.Kasten.Local.GitAnnexEnabled {
		err = p.files.Unlock()

		if err != nil {
			return
		}

		defer p.files.Lock()
	}

	cmd := util.ExecCommand(
		"open",
		[]string{"-W"},
		p.files.Files,
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
