package commands

import (
	"flag"
	"strings"

	"github.com/friedenberg/z/commands/options"
	"github.com/friedenberg/z/lib"
	"github.com/friedenberg/z/lib/feeder"
	"github.com/friedenberg/z/lib/pipeline"
	"github.com/friedenberg/z/lib/zettel"
	"github.com/friedenberg/z/lib/zettel/modifier"
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

	return func(u *lib.Umwelt) (err error) {
		p := pipeline.Pipeline{
			Feeder: feeder.MakeStringSlice(f.Args()),
			Reader: kind,
			Modifier: modifier.Chain(
				modifier.Make(
					func(i int, z *zettel.Zettel) (err error) {
						tags := strings.Split(tagString, " ")
						z.Metadata.AddStringTags(tags...)

						if description != "" {
							z.Metadata.SetDescription(description)
						}

						err = z.Write(nil)

						return
					},
				),
				lib.MakeTransactionAction(u.Transaction, lib.TransactionActionAdded),
			),
		}

		p.Run(u)

		err = u.Kasten.CommitTransaction(u)

		if err != nil {
			return
		}

		//this must come after the transaction is run, as this may be changed by the
		//transaction
		toAction := make([]string, 0, u.Transaction.Len())
		toAction = append(toAction, u.Transaction.ZettelsForActions(lib.TransactionActionAdded).Paths()...)
		toAction = append(toAction, u.Transaction.ZettelsForActions(lib.TransactionActionModified).Paths()...)

		u.Transaction = lib.MakeTransaction()

		p = pipeline.Pipeline{
			Feeder: feeder.MakeStringSlice(toAction),
			Modifier: modifier.Chain(
				&lib.ActionModifier{
					Umwelt:  u,
					Actions: editActions,
				},
				lib.MakeTransactionAction(u.Transaction, lib.TransactionActionModified),
			),
		}

		p.Run(u)

		return
	}
}
