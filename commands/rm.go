package commands

import (
	"flag"
	"os"

	"github.com/friedenberg/z/lib"
	"github.com/friedenberg/z/lib/pipeline/printer"
)

func init() {
	makeAndRegisterCommand(
		"rm",
		GetSubcommandRm,
	)
}

func GetSubcommandRm(f *flag.FlagSet) lib.Transactor {
	return func(u lib.Umwelt) (err error) {
		processor := MakeProcessor(
			u,
			f.Args(),
			&printer.NullZettelPrinter{},
		)

		processor.actioner = func(i int, z *lib.Zettel) (shouldPrint bool, actionErr error) {
			shouldPrint = true
			actionErr = os.Remove(z.Path)

			if actionErr != nil {
				return
			}

			if f, ok := z.Note.Metadata.LocalFile(); ok {
				actionErr = os.Remove(f.FilePath(u.BasePath))
			}

			u.Del.PrintZettel(0, z, actionErr)

			return
		}

		err = processor.Run()

		return
	}
}
