package commands

import (
	"flag"
	"path/filepath"

	"github.com/friedenberg/z/lib"
)

func GetSubcommandPrint(f *flag.FlagSet) CommandRunFunc {
	return func(e Env) (err error) {
		putter := MakePutter()
		glob := filepath.Join(e.ZettelPath, "*.md")
		processor, err := MakeProcessor(
			glob,
			func(z *lib.Zettel) {
				z.GenerateAlfredItemData()
			},
			putter,
		)

		if err != nil {
			//todo
		}

		err = processor.Run()
		return
	}
}
