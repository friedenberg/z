package commands

import (
	"flag"
	"strings"

	"github.com/friedenberg/z/commands/options"
	"github.com/friedenberg/z/lib"
	"github.com/friedenberg/z/lib/pipeline/printer"
	"github.com/friedenberg/z/util"
)

func init() {
	makeAndRegisterCommand(
		"add",
		GetSubcommandAdd,
	)
}

func GetSubcommandAdd(f *flag.FlagSet) lib.Transactor {
	var tagString string
	var description string
	var kind attachmentKind
	editActions := options.Actions(options.ActionEdit)

	f.Var(&editActions, "actions", "action to perform for the matched zettels")

	f.StringVar(&tagString, "tags", "", "parse the passed-in string as the metadata.")
	f.StringVar(&description, "description", "", "use this string as the zettel description")
	f.Var(&kind, "kind", "treat the positional arguments as this kind.")

	return func(u lib.Umwelt) (err error) {
		pr := &printer.MultiplexingZettelPrinter{
			Printer: &printer.ActionZettelPrinter{
				Umwelt: u,
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

			pr.PrintZettel(i, z, err)
			u.Add.PrintZettel(i, z, err)

			return
		}

		par := util.Parallelizer{
			Args:    f.Args(),
			Printer: pr,
		}

		par.Run(iter, errIterartion(pr.Printer))

		added := u.Added().Paths()

		err = u.RunTransaction(nil)

		if err != nil {
			return
		}

		u.Transaction = lib.MakeTransaction()

		action(u, added, nil, editActions)

		return
	}
}
