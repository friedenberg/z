package commands

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/friedenberg/z/lib"
)

func GetSubcommandClean(f *flag.FlagSet) CommandRunFunc {
	return func(e Env) (err error) {
		glob := filepath.Join(e.ZettelPath, "*.md")
		files, err := filepath.Glob(glob)

		processor := MakeProcessor(
			e,
			files,
			&NullPutter{Channel: make(PutterChannel)},
		)

		processor.hydrateAction = func(i int, z *lib.Zettel) error { return cleanZettel(z, false) }

		err = processor.Run()

		return
	}
}

func cleanZettel(z *lib.Zettel, dryRun bool) (err error) {
	didPrintPath := false
	printPathIfNecessary := func() {
		if !didPrintPath {
			fmt.Println(z.Path + ":")
		}

		didPrintPath = true
	}

	cleanActions := z.GetCleanActions()

	for n, a := range cleanActions {
		applicable := a.Check()

		if !applicable {
			continue
		}

		printPathIfNecessary()

		fmt.Fprintf(os.Stderr, "\t%s: yes\n", n)

		if !dryRun {
			a.Perform()
		}
	}

	return
}
