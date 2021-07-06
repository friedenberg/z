package commands

import (
	"flag"

	"github.com/friedenberg/z/lib"
)

func GetSubcommandOpen(f *flag.FlagSet) CommandRunFunc {
	var shouldEdit bool
	var shouldOpen bool

	f.BoolVar(&shouldEdit, "edit", true, "")
	f.BoolVar(&shouldOpen, "action", false, "")

	return func(e Env) (err error) {
		processor := MakeProcessor(
			e,
			f.Args(),
			&nullZettelPrinter{},
		)

		processor.actioner = func(i int, z *lib.Zettel) (err error) {
			if shouldOpen {
				z.Open(e.ZettelPath)

				if err != nil {
					return err
				}
			}

			if shouldEdit {
				err = z.Edit(e.ZettelPath)
			}

			return
		}

		err = processor.Run()

		return
	}
}
