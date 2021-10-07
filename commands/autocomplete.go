package commands

import (
	"flag"

	"github.com/friedenberg/z/lib"
)

func GetSubcommandAutocomplete(f *flag.FlagSet) lib.Transactor {
	return func(u lib.Umwelt, t lib.Transaction) error {
		return nil
	}
}
