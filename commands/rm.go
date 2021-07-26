package commands

import (
	"flag"
	"os"

	"github.com/friedenberg/z/lib"
)

func GetSubcommandRm(f *flag.FlagSet) CommandRunFunc {
	return func(e *lib.Env) (err error) {
		processor := MakeProcessor(
			e,
			f.Args(),
			&nullZettelPrinter{},
		)

		processor.actioner = func(i int, z *lib.Zettel) (actionErr error) {
			actionErr = os.Remove(z.Path)

			if actionErr != nil {
				return
			}

			if z.HasFile() {
				actionErr = os.Remove(z.FilePath())
			}

			return
		}

		err = processor.Run()

		return
	}
}
