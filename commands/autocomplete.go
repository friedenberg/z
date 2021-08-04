package commands

import (
	"flag"

	"github.com/friedenberg/z/lib"
)

func GetSubcommandAutocomplete(f *flag.FlagSet) CommandRunFunc {
	return func(e *lib.Kasten) error {
		return nil
	}
}
