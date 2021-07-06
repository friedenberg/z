package commands

import (
	"flag"

	"github.com/friedenberg/z/lib"
	"github.com/friedenberg/z/util"
)

func GetSubcommandClean(f *flag.FlagSet) CommandRunFunc {
	isDryRun := false

	f.BoolVar(&isDryRun, "dry-run", false, "")

	return func(e *lib.Env) (err error) {
		processor := MakeProcessor(
			e,
			f.Args(),
			&nullZettelPrinter{},
		)

		processor.actioner = cleanZettelFunc(isDryRun)

		err = processor.Run()

		return
	}
}

func cleanZettelFunc(dryRun bool) ActionFunc {
	return func(i int, z *lib.Zettel) (err error) {
		didPrintPath := false
		printPathIfNecessary := func() {
			if !didPrintPath {
				util.StdPrinterErr(z.Path + ":")
			}

			didPrintPath = true
		}

		cleanActions := lib.GetCleanActions()

		for n, a := range cleanActions {
			applicable := a.Check(z)

			if !applicable {
				continue
			}

			printPathIfNecessary()

			util.StdPrinterErrf("\t%s: yes\n", n)

			if !dryRun {
				z.ReadMetadataAndBody()
				a.Perform(z)
			}
		}

		return
	}
}
