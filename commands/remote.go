package commands

import (
	"flag"

	"github.com/friedenberg/z/commands/options"
	"github.com/friedenberg/z/lib"
	"github.com/friedenberg/z/lib/pipeline"
	"github.com/friedenberg/z/lib/pipeline/printer"
	"github.com/friedenberg/z/util"
)

func init() {
	makeAndRegisterCommand(
		"remote",
		GetSubcommandRemote,
	)
}

func GetSubcommandRemote(f *flag.FlagSet) lib.Transactor {
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

		remotePath := args[0]

		args = args[1:]

		////TODO-P3 validation
		//query := args[2]

		fp := pipeline.FilterPrinter{
			Printer: &printer.RemotePrinter{
				Umwelt:     u,
				Command:    command,
				RemotePath: remotePath,
			},
			Filter: func(_ int, z *lib.Zettel) bool {
				return z.Note.Metadata.HasFile()
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
