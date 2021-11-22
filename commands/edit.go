package commands

import (
	"flag"

	"github.com/friedenberg/z/commands/options"
	"github.com/friedenberg/z/lib"
	"github.com/friedenberg/z/lib/pipeline"
	"github.com/friedenberg/z/lib/zettel/filter"
	"github.com/friedenberg/z/lib/zettel/modifier"
	"golang.org/x/xerrors"
)

func init() {
	makeAndRegisterCommand(
		"edit",
		GetSubcommandEdit,
	)
}

func GetSubcommandEdit(f *flag.FlagSet) lib.Transactor {
	var query string
	var editAll bool
	editActions := options.Actions(options.ActionEdit)

	f.BoolVar(&editAll, "all", false, "edit all zettels if no arguments are passed in")
	f.StringVar(&query, "query", "", "zettel-spec string to determine which zettels to open or edit")
	f.Var(&editActions, "actions", "action to perform for the matched zettels")

	return func(u *lib.Umwelt) (err error) {
		args := f.Args()

		if len(args) == 0 {
			if editAll {
				err = xerrors.Errorf("refusing to edit all zettels unless '-all' flag is set")
				return
			} else {
				args = u.GetAll()
			}
		}

		p := pipeline.Pipeline{
			Arguments: args,
			Filter:    filter.Tag(query),
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
