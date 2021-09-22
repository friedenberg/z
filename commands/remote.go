package commands

import (
	"flag"

	"github.com/friedenberg/z/commands/options"
	"github.com/friedenberg/z/commands/printer"
	"github.com/friedenberg/z/lib"
	"github.com/friedenberg/z/lib/kasten"
	"github.com/friedenberg/z/lib/pipeline"
	"github.com/friedenberg/z/util"
	"golang.org/x/xerrors"
)

func GetSubcommandRemote(f *flag.FlagSet) CommandRunFunc {
	var query string
	f.StringVar(&query, "query", "", "zettel-spec")

	return func(u lib.Umwelt) (err error) {
		args := f.Args()

		var command options.RemoteCommand

		err = command.Set(args[0])

		if err != nil {
			return
		}

		args = args[1:]

		var remote kasten.RemoteImplementation
		var ok bool

		if remote, ok = u.RemoteKasten[args[0]]; !ok {
			err = xerrors.Errorf("invalid remote kasten: '%s'", args[1])
			return
		}

		args = args[1:]

		////TODO validation
		//query := args[2]

		fp := pipeline.FilterPrinter{
			Printer: &printer.RemotePrinter{
				Umwelt:  u,
				Command: command,
				Remote:  remote,
			},
		}

		var iter util.ParallelizerIterFunc

		if u.Config.UseIndexCache {
			if len(args) == 0 {
				args = u.GetAll()
			}

			iter = cachedIteration(u, query, fp)
		} else {
			if len(args) == 0 {
				args, err = u.FilesAndGit().GetAll()

				if err != nil {
					return
				}
			}

			iter = filesystemIteration(u, query, fp)
		}

		par := util.Parallelizer{Args: args}
		fp.Printer.Begin()
		defer fp.Printer.End()
		par.Run(iter, errIterartion(fp.Printer))

		return
	}
}