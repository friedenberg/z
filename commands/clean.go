package commands

import (
	"flag"
	"fmt"
	"os"
	"path"

	"github.com/friedenberg/z/lib"
)

func GetSubcommandClean(f *flag.FlagSet) CommandRunFunc {
	return func(e Env) (err error) {
		glob := path.Join(e.ZettelPath, "*.md")
		processor, err := MakeProcessor(
			glob,
			func(z *lib.Zettel) { cleanZettel(z, false) },
			&NullPutter{Channel: make(PutterChannel)},
		)

		if err != nil {
			return
		}

		err = processor.Run()

		return
	}
}

func cleanZettel(z *lib.Zettel, dryRun bool) {
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
}
