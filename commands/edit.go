package commands

import (
	"flag"

	"github.com/friedenberg/z/commands/printer"
	"github.com/friedenberg/z/lib"
)

func GetSubcommandEdit(f *flag.FlagSet) CommandRunFunc {
	var query string
	editActions := printer.Actions(printer.ActionEdit)

	f.StringVar(&query, "query", "", "zettel-spec string to determine which zettels to open or edit")
	f.Var(&editActions, "actions", "action to perform for the matched zettels")

	return func(e *lib.Kasten) (err error) {
		processor := MakeProcessor(
			e,
			f.Args(),
			&printer.MultiplexingZettelPrinter{
				Printer: &printer.ActionZettelPrinter{
					Kasten:  e,
					Actions: editActions,
				},
			},
		)

		processor.actioner = func(i int, z *lib.Zettel) (shouldPrint bool, err error) {
			shouldPrint = doesZettelMatchQuery(z, query)
			return
		}

		err = processor.Run()

		return
	}
}

//TODO refactor
func doesZettelMatchQuery(z *lib.Zettel, q string) bool {
	if q == "" {
		return true
	}

	if z.IndexData.File == q {
		return true
	}

	if z.IndexData.Url == q {
		return true
	}

	for _, t := range z.IndexData.ExpandedTags {
		if t == q {
			return true
		}
	}

	return false
}
