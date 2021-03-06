package commands

import (
	"flag"

	"github.com/friedenberg/z/lib"
	"github.com/friedenberg/z/lib/pipeline"
	"github.com/friedenberg/z/lib/zettel/reader"
)

func init() {
	makeAndRegisterCommand(
		"index",
		GetSubcommandIndex,
	)
}

func GetSubcommandIndex(f *flag.FlagSet) lib.Transactor {
	return func(u *lib.Umwelt) (err error) {
		u.ShouldSkipCommit = true
		u.Index = lib.MakeIndex()

		args := u.Kasten.GetAll()

		p := pipeline.Pipeline{
			Feeder:   args,
			Reader:   reader.FromFile(true),
			Modifier: lib.MakeTransactionAction(u.Transaction, lib.TransactionActionAdded),
		}

		p.Run(u)

		return
	}
}
