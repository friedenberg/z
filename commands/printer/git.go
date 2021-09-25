package printer

import (
	"sync"

	"github.com/friedenberg/z/lib"
	"github.com/friedenberg/z/util"
)

type GitPrinter struct {
	*sync.Mutex
	Umwelt           lib.Umwelt
	GitCommitMessage string
	shouldCommit     bool
	files            []string
}

func (p *GitPrinter) Begin() {
	if !p.Umwelt.Kasten.Local.GitEnabled {
		return
	}
}

func (p *GitPrinter) SetShouldCommit() {
	p.Lock()
	defer p.Unlock()
	p.shouldCommit = true
}

func (p *GitPrinter) PrintZettel(i int, z *lib.Zettel, errIn error) {
	if !p.Umwelt.Kasten.Local.GitEnabled {
		return
	}

	if errIn != nil {
		util.StdPrinterErr(errIn)
		return
	}

	p.files = append(p.files, z.Path)

	if z.HasFile() {
		p.files = append(p.files, z.FilePath())
	}
}

func (p *GitPrinter) End() {
	if !p.Umwelt.Kasten.Local.GitEnabled {
		return
	}

	var err error

	if p.shouldCommit && len(p.files) > 0 {
		git := util.GitFilesToCommit{
			Git: util.Git{
				Path: p.Umwelt.Kasten.Local.BasePath,
			},
			Files: p.files,
		}

		err = git.AddAndCommit(p.GitCommitMessage)

		if err != nil {
			util.StdPrinterErr(err)
		}
	}
}
