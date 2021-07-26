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
		cleanActions := lib.GetCleanActions()
		shouldWrite := false

		for n, a := range cleanActions {
			applicable := a.Check(z)

			if !applicable {
				continue
			}

			util.StdPrinterErrf("%s: %s\n", z.Path, n)

			if !dryRun {
				z.ReadMetadataAndBody()
				var newShouldWrite bool
				newShouldWrite, err = a.Perform(z)
				shouldWrite = shouldWrite || newShouldWrite

				if err != nil {
					return
				}
			}
		}

		if shouldWrite {
			util.OpenFilesGuardInstance.Lock()
			defer util.OpenFilesGuardInstance.Unlock()
			err = z.Write(nil)
		}

		return
	}
}
