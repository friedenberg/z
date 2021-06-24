package commands

import (
	"flag"
	"os"

	"github.com/friedenberg/z/lib"
)

func GetSubcommandRm(f *flag.FlagSet) CommandRunFunc {
	return func(e Env) (err error) {
		path := f.Arg(0)

		if path == "" {
			//TODO
		}

		z := &lib.Zettel{}
		z.HydrateFromFilePath(path)

		err = os.Remove(z.Path)

		if err != nil {
			return
		}

		if z.Metadata.Kind == "file" {
			err = os.Remove(z.Metadata.File)
		}

		return
	}
}
