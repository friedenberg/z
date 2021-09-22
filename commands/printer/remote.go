package printer

import (
	"path"

	"github.com/friedenberg/z/commands/options"
	"github.com/friedenberg/z/lib"
	"github.com/friedenberg/z/lib/kasten"
	"github.com/friedenberg/z/util"
	"golang.org/x/xerrors"
)

type RemotePrinter struct {
	Umwelt  lib.Umwelt
	Command options.RemoteCommand
	Remote  kasten.RemoteImplementation
	zettels []*lib.Zettel
}

func (p *RemotePrinter) Begin() {
	p.zettels = make([]*lib.Zettel, 0)
}

func (p *RemotePrinter) PrintZettel(i int, z *lib.Zettel, errIn error) {
	if errIn != nil {
		util.StdPrinterError(errIn)
		return
	}

	var err error

	switch p.Command {
	case options.RemoteCommandPull:

		util.StdPrinterErrf("%s: copy start\n", z.FilePath())

		err = p.Remote.CopyFileFrom(
			z.FilePath(),
			kasten.RemoteFileHandle{z.Id, path.Ext(z.FilePath())},
		)

	case options.RemoteCommandPush:

		util.StdPrinterErrf("%s: copy start\n", z.FilePath())

		err = p.Remote.CopyFileTo(
			z.FilePath(),
			kasten.RemoteFileHandle{z.Id, path.Ext(z.FilePath())},
		)

	default:
		panic(xerrors.Errorf("unsupported remote command: '%s'", p.Command))
	}

	if err == nil {
		util.StdPrinterErrf("%s: copy end\n", z.FilePath())
		p.zettels = append(p.zettels, z)
	} else {
		util.StdPrinterError(xerrors.Errorf("%s: copy end: %w", z.FilePath(), err))
	}
}

func (p *RemotePrinter) End() {
	//TODO record location of files in remote in zettels
	// p.Umwelt.LocalKasten
}
