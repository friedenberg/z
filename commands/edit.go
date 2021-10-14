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
		"edit",
		GetSubcommandEdit,
	)
}

func GetSubcommandEdit(f *flag.FlagSet) lib.Transactor {
	var query string
	editActions := options.Actions(options.ActionEdit)

	f.StringVar(&query, "query", "", "zettel-spec string to determine which zettels to open or edit")
	f.Var(&editActions, "actions", "action to perform for the matched zettels")

	return func(u lib.Umwelt) (err error) {
		fp := pipeline.FilterPrinter{
			Filter: MatchQuery(query),
			Printer: &printer.MultiplexingZettelPrinter{
				Printer: &printer.ActionZettelPrinter{
					Umwelt:  u,
					Actions: editActions,
				},
			},
		}

		args := f.Args()
		var iter util.ParallelizerIterFunc

		if len(args) == 0 {
			args = u.GetAll()
		}

		iter = cachedIteration(u, query, fp)

		par := util.Parallelizer{Args: args}
		fp.Printer.Begin()
		defer fp.Printer.End()
		par.Run(
			iter,
			errIterartion(fp.Printer),
		)

		return
	}
}
