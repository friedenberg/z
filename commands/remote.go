package commands

import (
	"flag"

	"github.com/friedenberg/z/commands/options"
	"github.com/friedenberg/z/lib"
	"github.com/friedenberg/z/lib/kasten"
	"github.com/friedenberg/z/lib/pipeline"
	"github.com/friedenberg/z/lib/pipeline/printer"
	"github.com/friedenberg/z/util"
	"golang.org/x/xerrors"
)

func GetSubcommandRemote(f *flag.FlagSet) lib.Transactor {
	var query string
	f.StringVar(&query, "query", "", "zettel-spec")

	return func(u lib.Umwelt, t lib.Transaction) (err error) {
		args := f.Args()

		var command options.RemoteCommand

		err = command.Set(args[0])

		if err != nil {
			return
		}

		args = args[1:]

		var remote kasten.RemoteImplementation
		var ok bool

		if remote, ok = u.Kasten.Remotes[args[0]]; !ok {
			err = xerrors.Errorf("invalid remote kasten: '%s'", args[1])
			return
		}

		args = args[1:]

		////TODO validation
		//query := args[2]

		fp := pipeline.FilterPrinter{
			Printer: &printer.RemotePrinter{
				Umwelt:      u,
				Transaction: t,
				Command:     command,
				Remote:      remote,
			},
			Filter: func(_ int, z *lib.Zettel) bool {
				return z.HasFile()
			},
		}

		var iter util.ParallelizerIterFunc

		if len(args) == 0 {
			args = u.GetAll()
		}

		iter = cachedIteration(u, query, fp)

		par := util.Parallelizer{Args: args}
		fp.Printer.Begin()
		defer fp.Printer.End()
		par.Run(iter, errIterartion(fp.Printer))

		return
	}
}
