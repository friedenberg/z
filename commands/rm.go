package commands

import (
	"flag"
	"os"

	"github.com/friedenberg/z/commands/printer"
	"github.com/friedenberg/z/lib"
)

func GetSubcommandRm(f *flag.FlagSet) CommandRunFunc {
	return func(e *lib.Kasten) (err error) {
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
