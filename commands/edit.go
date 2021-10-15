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
		args := f.Args()

		if len(args) == 0 {
			//TODO-P3 does it make sense to edit all zettels on no args?
			args = u.GetAll()
		}

		err = action(
			u,
			args,
			pipeline.MatchQuery(query),
			editActions,
		)

		return
	}
}

func action(u lib.Umwelt, args []string, f pipeline.Filter, a options.Actions) (err error) {
	fp := pipeline.FilterPrinter{
		Filter: f,
		Printer: &printer.MultiplexingZettelPrinter{
			Printer: &printer.ActionZettelPrinter{
				Umwelt:  u,
				Actions: a,
			},
		},
	}

	iter := cachedIteration(u, fp)

	par := util.Parallelizer{
		Args:    args,
		Printer: fp.Printer,
	}

	par.Run(
		iter,
		errIterartion(fp.Printer),
	)

	return
}
