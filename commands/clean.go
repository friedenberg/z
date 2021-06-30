package commands

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/friedenberg/z/lib"
)

func GetSubcommandClean(f *flag.FlagSet) CommandRunFunc {
	isDryRun := false

	f.BoolVar(&isDryRun, "dry-run", false, "")

	return func(e Env) (err error) {
		glob := filepath.Join(e.ZettelPath, "*.md")
		files, err := filepath.Glob(glob)

		processor := MakeProcessor(
			e,
			files,
			&NullPutter{Channel: make(PutterChannel)},
		)

		processor.parallelAction = cleanZettelFunc(isDryRun)

		err = processor.Run()

		return
	}
}

func cleanZettelFunc(dryRun bool) ProcessorAction {
	return func(i int, z *lib.Zettel) (err error) {
		didPrintPath := false
		printPathIfNecessary := func() {
			if !didPrintPath {
				fmt.Println(z.Path + ":")
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

			fmt.Fprintf(os.Stderr, "\t%s: yes\n", n)

			if !dryRun {
				a.Perform(z)
			}
		}

		return
	}
}
