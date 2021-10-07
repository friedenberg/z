package commands

import (
	"flag"

	"github.com/friedenberg/z/lib"
)

func GetSubcommandIndex(f *flag.FlagSet) lib.Transactor {
	return func(u lib.Umwelt, t lib.Transaction) (err error) {
		u.Index = lib.MakeIndex()
		err = hydrateIndex(u)

		if err != nil {
			return
		}

		err = u.CacheIndex()

		if err != nil {
			return
		}

		return
	}
}
