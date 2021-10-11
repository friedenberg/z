package commands

import (
	"flag"
	"strings"

	"github.com/friedenberg/z/commands/options"
	"github.com/friedenberg/z/lib"
	"github.com/friedenberg/z/lib/pipeline/printer"
	"github.com/friedenberg/z/util"
)

func GetSubcommandAdd(f *flag.FlagSet) lib.Transactor {
	var tagString string
	var description string
	var kind attachmentKind
	editActions := options.Actions(options.ActionEdit)

	f.Var(&editActions, "actions", "action to perform for the matched zettels")

	f.StringVar(&tagString, "tags", "", "parse the passed-in string as the metadata.")
	f.StringVar(&description, "description", "", "use this string as the zettel description")
	f.Var(&kind, "kind", "treat the positional arguments as this kind.")

	return func(u lib.Umwelt, t *lib.Transaction) (err error) {
		pr := &printer.MultiplexingZettelPrinter{
			Printer: &printer.ActionZettelPrinter{
				Umwelt:      u,
				Actions:     editActions,
				Transaction: t,
			},
		}

		iter := func(i int, a string) (err error) {
			z, err := kind.hydrator(u, a)

			if err != nil {
				return
			}

			tags := strings.Split(tagString, " ")
			z.Metadata.AddStringTags(tags...)

			if description != "" {
				z.Metadata.SetDescription(description)
			}

			t.Add.PrintZettel(i, z, err)

			return
		}

		par := util.Parallelizer{Args: f.Args()}
		pr.Printer.Begin()
		defer pr.Printer.End()
		par.Run(iter, errIterartion(pr.Printer))

		return
	}
}
