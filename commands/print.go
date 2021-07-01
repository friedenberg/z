package commands

import (
	"flag"
	"path/filepath"

	"github.com/friedenberg/z/lib"
)

func GetSubcommandPrint(f *flag.FlagSet) CommandRunFunc {
	return func(e Env) (err error) {
		glob := filepath.Join(e.ZettelPath, "*.md")
		files, err := filepath.Glob(glob)

		if err != nil {
			return
		}

		processor := MakeProcessor(
			e,
			files,
			MakePutter(),
		)

		processor.actioner = func(i int, z *lib.Zettel) error {
			return z.GenerateAlfredItemData()
		}

		err = processor.Run()

		return
	}
}
