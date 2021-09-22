package commands

import (
	"flag"

	"github.com/friedenberg/z/commands/options"
	"github.com/friedenberg/z/commands/printer"
	"github.com/friedenberg/z/lib"
	"github.com/friedenberg/z/lib/pipeline"
	"github.com/friedenberg/z/util"
)

func GetSubcommandEdit(f *flag.FlagSet) CommandRunFunc {
	var query string
	editActions := options.Actions(options.ActionEdit)

	f.StringVar(&query, "query", "", "zettel-spec string to determine which zettels to open or edit")
	f.Var(&editActions, "actions", "action to perform for the matched zettels")

	return func(e lib.Umwelt) (err error) {
		fp := pipeline.FilterPrinter{
			Filter: MatchQuery(query),
			Printer: &printer.MultiplexingZettelPrinter{
				Printer: &printer.ActionZettelPrinter{
					Umwelt:  e,
					Actions: editActions,
				},
			},
		}

		args := f.Args()
		var iter util.ParallelizerIterFunc

		if e.Config.UseIndexCache {
			if len(args) == 0 {
				args = e.GetAll()
			}

			iter = cachedIteration(e, query, fp)
		} else {
			if len(args) == 0 {
				args, err = e.FilesAndGit().GetAll()

				if err != nil {
					return
				}
			}

			iter = filesystemIteration(e, query, fp)
		}

		par := util.Parallelizer{Args: args}
		fp.Printer.Begin()
		defer fp.Printer.End()
		par.Run(iter, errIterartion(fp.Printer))

		return
	}
}
