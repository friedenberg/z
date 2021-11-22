package lib

import (
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/friedenberg/z/commands/options"
	"github.com/friedenberg/z/lib/zettel"
	"github.com/friedenberg/z/lib/zettel/metadata"
	"github.com/friedenberg/z/util"
	"github.com/friedenberg/z/util/stdprinter"
	"golang.org/x/xerrors"
)

type ActionModifier struct {
	Umwelt      *Umwelt
	Actions     options.Actions
	zettels     []*zettel.Zettel
	zettelFiles []string
	files       []string
	urls        []metadata.Url
}

func (p *ActionModifier) ModifyZettel(i int, z *zettel.Zettel) (err error) {
	p.zettels = append(p.zettels, z)
	p.zettelFiles = append(p.zettelFiles, z.Path)

	if f, ok := z.Note.Metadata.LocalFile(); ok {
		p.files = append(p.files, f.FilePath(p.Umwelt.BasePath))
	}

	if u, ok := z.Note.Metadata.Url(); ok {
		p.urls = append(p.urls, u)
	}

	if p.Actions&options.ActionPrintZettelPath != 0 {
		_, err := io.WriteString(os.Stdout, fmt.Sprintf("%s\n", z.Path))
		stdprinter.PanicIfError(err)
	}

	return
}

func (p *ActionModifier) End(_ io.Writer) {
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

	if err != nil {
		stdprinter.Err(err)
	}
}

func (p *ActionModifier) openZettels() (err error) {
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

func (p *ActionModifier) openFiles() (err error) {
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

func (p *ActionModifier) openUrls() (err error) {
	if len(p.urls) == 0 {
		return
	}

	fileUs := make([]string, 0)
	normalUs := make([]string, 0)

	for _, u := range p.urls {
		cs := u.CorrectedString()

		if u.Scheme == "file" {
			fileUs = append(fileUs, cs)
		} else {
			normalUs = append(normalUs, cs)
		}
	}

	err = p.openNormalUrls(normalUs)

	if err != nil {
		return
	}

	err = p.openFileUrls(fileUs)

	if err != nil {
		return
	}

	return
}

func (p *ActionModifier) openNormalUrls(us []string) (err error) {
	args := []string{
		"-na",
		"Google Chrome",
		"--args",
		"--new-window",
	}

	cmd := util.ExecCommand(
		"open",
		args,
		us,
	)

	output, err := cmd.CombinedOutput()

	if err != nil {
		err = xerrors.Errorf("opening urls ('%q'): %s", p.urls, output)
		return
	}

	return
}

func (p *ActionModifier) openFileUrls(us []string) (err error) {
	cmd := util.ExecCommand(
		"open",
		us,
	)

	output, err := cmd.CombinedOutput()

	if err != nil {
		err = xerrors.Errorf("opening file urls ('%q'): %s", p.urls, output)
		return
	}

	return
}
