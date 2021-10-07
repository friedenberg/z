package commands

import (
	"flag"
	"os"

	"github.com/friedenberg/z/lib"
	"github.com/friedenberg/z/lib/pipeline/printer"
)

func GetSubcommandRm(f *flag.FlagSet) CommandRunFunc {
	return func(e lib.Umwelt) (err error) {
		processor := MakeProcessor(
			e,
			f.Args(),
			&printer.NullZettelPrinter{},
		)

		processor.actioner = func(i int, z *lib.Zettel) (shouldPrint bool, actionErr error) {
			shouldPrint = true
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
