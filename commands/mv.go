package commands

import (
	"flag"

	"github.com/friedenberg/z/lib"
)

func GetSubcommandMv(f *flag.FlagSet) lib.Transactor {
	isDryRun := false

	f.BoolVar(&isDryRun, "dry-run", false, "")

	return func(u lib.Umwelt, t *lib.Transaction) (err error) {
		//TODO

		return
	}
}
