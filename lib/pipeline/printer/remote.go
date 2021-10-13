package printer

import (
	"github.com/friedenberg/z/commands/options"
	"github.com/friedenberg/z/lib"
	"github.com/friedenberg/z/lib/kasten"
	"github.com/friedenberg/z/util"
	"golang.org/x/xerrors"
)

type RemotePrinter struct {
	Umwelt       lib.Umwelt
	Command      options.RemoteCommand
	Remote       kasten.RemoteImplementation
	zettels      []*lib.Zettel
	rsyncPrinter *Rsync
}

func (p *RemotePrinter) Begin() {
	p.rsyncPrinter = &Rsync{}

	switch p.Command {
	case options.RemoteCommandPull:
		p.rsyncPrinter.Src = p.Remote.(*kasten.Files).BasePath
		p.rsyncPrinter.Dst = p.Umwelt.BasePath

	case options.RemoteCommandPush:
		p.rsyncPrinter.Src = p.Umwelt.BasePath
		p.rsyncPrinter.Dst = p.Remote.(*kasten.Files).BasePath

	default:
		panic(xerrors.Errorf("unsupported remote command: '%s'", p.Command))
	}

	p.zettels = make([]*lib.Zettel, 0)
	p.rsyncPrinter.Begin()
}

func (p *RemotePrinter) PrintZettel(i int, z *lib.Zettel, errIn error) {
	if errIn != nil {
		util.StdPrinterError(errIn)
		return
	}

	fd, ok := z.Note.Metadata.LocalFile()

	if !ok {
		//TODO-P4 decide whether to skip or to error
		util.StdPrinterError(xerrors.Errorf("zettel ('%s') has no file descriptors", z.Id))
		return
	}

	p.rsyncPrinter.File(fd.FileName())
	z.Note.Metadata.AddFile(fd)
	p.Umwelt.Mod.PrintZettel(i, z, errIn)
}

func (p *RemotePrinter) End() {
	p.rsyncPrinter.End()
	//TODO-P0 record location of files in remote in zettels
	// p.Umwelt.LocalKasten
}
