package printer

import (
	"fmt"
	"os/exec"
	"sync"

	"github.com/friedenberg/z/lib"
	"github.com/friedenberg/z/util"
)

type ActionZettelPrinter struct {
	Umwelt      lib.Umwelt
	Actions     Actions
	zettels     []*lib.Zettel
	zettelFiles []string
	files       []string
	urls        []string
}

func (p *ActionZettelPrinter) Begin() {}

func (p *ActionZettelPrinter) PrintZettel(i int, z *lib.Zettel, errIn error) {
	if errIn != nil {
		util.StdPrinterErr(errIn)
		return
	}

	p.zettels = append(p.zettels, z)
	p.zettelFiles = append(p.zettelFiles, z.Path)

	if z.HasFile() {
		p.files = append(p.files, z.FilePath())
	}

	if z.HasUrl() {
		p.urls = append(p.urls, z.Metadata.Url)
	}

	if p.Actions&ActionPrintZettelPath != 0 {
		//TODO full path
		util.StdPrinterOut(z.Path)
	}
}

func (p *ActionZettelPrinter) End() {
	gitPrinter := &GitPrinter{
		Umwelt:           p.Umwelt,
		Mutex:            &sync.Mutex{},
		GitCommitMessage: "edit",
	}

	gitPrinter.Begin()
	defer gitPrinter.End()

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

			gitPrinter.SetShouldCommit()
		}()
	}

	if p.Actions&ActionEdit != 0 {
		runAction(p.openZettels)
	}

	if p.Actions&ActionOpenFile != 0 {
		runAction(p.openFiles)
	}

	if p.Actions&ActionOpenUrl != 0 {
		runAction(p.openUrls)
	}

	wg.Wait()

	for i, z := range p.zettels {
		gitPrinter.PrintZettel(i, z, nil)
	}

	if err != nil {
		util.StdPrinterErr(err)
	}
}

func (p *ActionZettelPrinter) openZettels() (err error) {
	if len(p.zettels) == 0 {
		return
	}

	ga := &util.GitAnnex{
		GitFilesToCommit: util.GitFilesToCommit{
			Git: util.Git{
				Path: p.Umwelt.FilesAndGit().BasePath,
			},
			Files: p.zettelFiles,
		},
	}

	err = ga.Unlock()

	if err != nil {
		return
	}

	defer ga.Lock()

	args := []string{"-f", "-p"}

	cmd := exec.Command(
		"mvim",
		append(args, p.zettelFiles...)...,
	)

	output, err := cmd.CombinedOutput()

	if err != nil {
		err = fmt.Errorf("opening zettels ('%q'): %s", p.zettels, output)
		return
	}

	return
}

func (p *ActionZettelPrinter) openFiles() (err error) {
	if len(p.files) == 0 {
		return
	}

	ga := &util.GitAnnex{
		GitFilesToCommit: util.GitFilesToCommit{
			Git: util.Git{
				Path: p.Umwelt.FilesAndGit().BasePath,
			},
			Files: p.files,
		},
	}

	err = ga.Unlock()

	if err != nil {
		return
	}

	defer ga.Lock()

	cmd := exec.Command(
		"open",
		append([]string{"-W"}, p.files...)...,
	)

	output, err := cmd.CombinedOutput()

	if err != nil {
		err = fmt.Errorf("opening files ('%q'): %s", p.files, output)
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

	cmd := exec.Command(
		"open",
		append(args, p.urls...)...,
	)

	output, err := cmd.CombinedOutput()

	if err != nil {
		err = fmt.Errorf("opening urls ('%q'): %s", p.urls, output)
		return
	}

	return
}
