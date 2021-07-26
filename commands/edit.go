package commands

import (
	"flag"

	"github.com/friedenberg/z/lib"
)

func GetSubcommandEdit(f *flag.FlagSet) CommandRunFunc {
	var shouldEdit bool
	var shouldOpen bool

	f.BoolVar(&shouldEdit, "edit", true, "")
	f.BoolVar(&shouldOpen, "open", false, "")

	return func(e *lib.Env) (err error) {
		processor := MakeProcessor(
			e,
			f.Args(),
			&nullZettelPrinter{},
		)

		processor.actioner = func(i int, z *lib.Zettel) (err error) {
			if shouldEdit {
				z.Edit()

				if err != nil {
					return err
				}
			}

			if shouldOpen {
				err = z.Open()
			}

			return
		}

		err = processor.Run()

		return
	}
}
