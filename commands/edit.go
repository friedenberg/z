package commands

import (
	"flag"

	"github.com/friedenberg/z/commands/printer"
	"github.com/friedenberg/z/lib"
	"github.com/friedenberg/z/util"
)

func GetSubcommandEdit(f *flag.FlagSet) CommandRunFunc {
	var query string
	editActions := printer.Actions(printer.ActionEdit)

	f.StringVar(&query, "query", "", "zettel-spec string to determine which zettels to open or edit")
	f.Var(&editActions, "actions", "action to perform for the matched zettels")

	return func(e lib.Umwelt) (err error) {
		processor := MakeProcessor(
			e,
			f.Args(),
			&printer.MultiplexingZettelPrinter{
				Printer: &printer.ActionZettelPrinter{
					Umwelt:  e,
					Actions: editActions,
				},
			},
		)

		processor.argNormalizer = func(_ int, p string) (normalizedArg string, err error) {
			b := util.BaseNameNoSuffix(p)
			p = b + ".md"
			normalizedArg, err = e.FilesAndGit().GetNormalizedPath(p)
			return
		}

		processor.actioner = func(i int, z *lib.Zettel) (shouldPrint bool, err error) {
			shouldPrint = doesZettelMatchQuery(z, query)
			return
		}

		err = processor.Run()

		return
	}
}
