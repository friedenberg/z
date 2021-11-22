package commands

import (
	"flag"

	"github.com/friedenberg/z/lib"
	"github.com/friedenberg/z/lib/pipeline"
	"github.com/friedenberg/z/lib/zettel"
	"github.com/friedenberg/z/lib/zettel/filter"
	"github.com/friedenberg/z/lib/zettel/modifier"
)

func init() {
	makeAndRegisterCommand(
		"mv",
		GetSubcommandMv,
	)
}

func GetSubcommandMv(f *flag.FlagSet) lib.Transactor {
	isDryRun := false

	f.BoolVar(&isDryRun, "dry-run", false, "")

	return func(u *lib.Umwelt) (err error) {
		t1 := f.Args()[0]
		t2 := f.Args()[1]

		p := pipeline.Pipeline{
			Arguments: u.GetAll(),
			Filter:    filter.String(t1),
			Modifier: modifier.Chain(
				modifier.Make(
					func(i int, z *zettel.Zettel) (err error) {
						z.Metadata.Delete(t1)
						z.Metadata.AddStringTags(t2)

						z.Write(nil)

						return
					},
				),
				lib.MakeTransactionAction(u.Transaction, lib.TransactionActionModified),
			),
		}

		p.Run(u)

		return
	}
}
