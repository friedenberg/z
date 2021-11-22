package commands

import (
	"flag"

	"github.com/friedenberg/z/lib"
	"github.com/friedenberg/z/lib/feeder"
	"github.com/friedenberg/z/lib/pipeline"
	"golang.org/x/xerrors"
)

func init() {
	makeAndRegisterCommand(
		"rm",
		GetSubcommandRm,
	)
}

func GetSubcommandRm(f *flag.FlagSet) lib.Transactor {
	return func(u *lib.Umwelt) (err error) {
		var args feeder.Feeder

		if len(f.Args()) == 0 {
			err = xerrors.Errorf("no zettels included for deletion")
			return
		} else {
			args = feeder.MakeStringSlice(f.Args())
		}

		p := pipeline.Pipeline{
			Feeder:   args,
			Modifier: lib.MakeTransactionAction(u.Transaction, lib.TransactionActionDeleted),
		}

		p.Run(u)

		return
	}
}
