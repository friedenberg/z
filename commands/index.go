package commands

import (
	"flag"

	"github.com/friedenberg/z/lib"
)

func GetSubcommandIndex(f *flag.FlagSet) CommandRunFunc {
	return func(e lib.Umwelt) (err error) {
		e.Index = lib.MakeIndex()
		err = hydrateIndex(e)

		if err != nil {
			return
		}

		if err != nil {
			return
		}

		err = e.CacheIndex()

		if err != nil {
			return
		}

		return
	}
}
