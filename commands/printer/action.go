package printer

import (
	"os/exec"

	"github.com/friedenberg/z/lib"
	"github.com/friedenberg/z/util"
)

type ActionZettelPrinter struct {
	Kasten                    *lib.Kasten
	ShouldEdit, ShouldOpen bool
	zettels                []*lib.Zettel
	files                  []string
	urls                   []string
}

func (p *ActionZettelPrinter) Begin() {}

func (p *ActionZettelPrinter) PrintZettel(i int, z *lib.Zettel, errIn error) {
	if errIn != nil {
		util.StdPrinterErr(errIn)
		return
	}

	p.zettels = append(p.zettels, z)

	if z.HasFile() {
		p.files = append(p.files, z.FilePath())
	}

	if z.HasUrl() {
		p.urls = append(p.urls, z.IndexData.Url)
	}
}

func (p *ActionZettelPrinter) End() {
	if p.ShouldEdit {
		p.openZettels()
	}

	if p.ShouldOpen {
		p.openFiles()
		p.openUrls()
	}
}

func (p *ActionZettelPrinter) openZettels() {
	if len(p.zettels) == 0 {
		return
	}

	zettelFiles := make([]string, len(p.zettels))

	for i, z := range p.zettels {
		zettelFiles[i] = z.Path
	}

	args := []string{"-p"}

	cmd := exec.Command(
		"mvim",
		append(args, zettelFiles...)...,
	)

	cmd.Run()
}

func (p *ActionZettelPrinter) openFiles() {
	if len(p.files) == 0 {
		return
	}

	cmd := exec.Command(
		"open",
		p.files...,
	)

	cmd.Run()
}

func (p *ActionZettelPrinter) openUrls() {
	if len(p.urls) == 0 {
		return
	}

	args := []string{
		"-na",
		"Google Chrome",
		"--args",
		"--new-window",
	}

	cmd := exec.Command(
		"open",
		append(args, p.urls...)...,
	)

	cmd.Run()
	//TODO return errors
}
