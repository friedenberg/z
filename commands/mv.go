package commands

import (
	"flag"

	"github.com/friedenberg/z/lib"
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
		//TODO

		return
	}
}
