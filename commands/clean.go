package commands

import (
	"flag"
	"sync"

	"github.com/friedenberg/z/lib"
	"github.com/friedenberg/z/lib/pipeline/printer"
	"github.com/friedenberg/z/util"
)

func GetSubcommandClean(f *flag.FlagSet) lib.Transactor {
	isDryRun := false

	f.BoolVar(&isDryRun, "dry-run", false, "")

	return func(u lib.Umwelt, t lib.Transaction) (err error) {
		processor := MakeProcessor(
			u,
			f.Args(),
			&printer.NullZettelPrinter{},
		)

		for n, a := range lib.GetCleanActions() {
			gitPrinter := &printer.GitPrinter{
				Mutex:            &sync.Mutex{},
				GitCommitMessage: n,
				Umwelt:           u,
			}

			gitPrinter.SetShouldCommit()

			processor.SetPrinter(gitPrinter)

			processor.actioner = cleanZettelFunc(isDryRun, n, a)

			err = processor.Run()

			if err != nil {
				return
			}
		}

		return
	}
}

func cleanZettelFunc(dryRun bool, name string, cleanAction lib.CleanAction) ActionFunc {
	return func(i int, z *lib.Zettel) (shouldPrint bool, err error) {
		shouldPrint = true
		shouldWrite := false

		applicable := cleanAction.Check(z)

		if !applicable {
			return
		}

		util.StdPrinterErrf("%s: %s\n", z.Path, name)

		if !dryRun {
			z.ReadMetadataAndBody()
			var newShouldWrite bool
			newShouldWrite, err = cleanAction.Perform(z)
			shouldWrite = shouldWrite || newShouldWrite

			if err != nil {
				return
			}
		}

		if shouldWrite {
			err = z.Write(nil)
		}

		return
	}
}
