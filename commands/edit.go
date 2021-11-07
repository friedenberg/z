package commands

import (
	"flag"

	"github.com/friedenberg/z/commands/options"
	"github.com/friedenberg/z/lib"
	"github.com/friedenberg/z/lib/pipeline"
	"github.com/friedenberg/z/lib/zettel/filter"
	"github.com/friedenberg/z/lib/zettel/modifier"
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

	return func(u *lib.Umwelt) (err error) {
		args := f.Args()

		if len(args) == 0 {
			//TODO-P3 does it make sense to edit all zettels on no args?
			args = u.GetAll()
		}

		p := pipeline.Pipeline{
			Arguments: args,
			Filter:    filter.MatchQuery(query),
			Modifier: modifier.Chain(
				&modifier.Action{
					Umwelt:  *u,
					Actions: editActions,
				},
				modifier.TransactionAction(u.Transaction, lib.TransactionActionModified),
			),
		}

		p.Run(u)

		return
	}
}
