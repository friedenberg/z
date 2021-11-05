package commands

import (
	"flag"

	"github.com/friedenberg/z/lib"
	"github.com/friedenberg/z/lib/pipeline"
	"github.com/friedenberg/z/lib/pipeline/modifier"
	"github.com/friedenberg/z/lib/pipeline/reader"
)

func init() {
	makeAndRegisterCommand(
		"index",
		GetSubcommandIndex,
	)
}

func GetSubcommandIndex(f *flag.FlagSet) lib.Transactor {
	return func(u lib.Umwelt) (err error) {
		u.ShouldSkipCommit = true
		u.Index = lib.MakeIndex()

		args, err := u.Kasten.GetAll()

		if err != nil {
			return
		}

		p := pipeline.Pipeline{
			Arguments: args,
			Reader:    reader.FromFile(true),
			Modifier:  modifier.TransactionAction(u.Transaction, lib.TransactionActionAdded),
		}

		p.Run(u)

		return
	}
}
