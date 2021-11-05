package commands

import (
	"flag"

	"github.com/friedenberg/z/lib"
	"github.com/friedenberg/z/lib/pipeline"
	"github.com/friedenberg/z/lib/pipeline/modifier"
	"golang.org/x/xerrors"
)

func init() {
	makeAndRegisterCommand(
		"rm",
		GetSubcommandRm,
	)
}

func GetSubcommandRm(f *flag.FlagSet) lib.Transactor {
	return func(u lib.Umwelt) (err error) {
		args := f.Args()

		if len(args) == 0 {
			err = xerrors.Errorf("no zettels included for deletion")
			return
		}

		p := pipeline.Pipeline{
			Arguments: args,
			Modifier:  modifier.TransactionAction(u.Transaction, lib.TransactionActionDeleted),
		}

		p.Run(u)

		return
	}
}
