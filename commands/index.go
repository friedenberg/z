package commands

import (
	"flag"

	"github.com/friedenberg/z/lib"
	"github.com/friedenberg/z/lib/pipeline"
	"github.com/friedenberg/z/util"
)

func GetSubcommandIndex(f *flag.FlagSet) lib.Transactor {
	return func(u lib.Umwelt, t lib.Transaction) (err error) {
		u.Index = lib.MakeIndex()

		args, err := u.FilesAndGit().GetAll()

		if err != nil {
			return
		}

		par := util.Parallelizer{Args: args}
		t.Add.Begin()
		defer t.Add.End()
		par.Run(
			func(i int, a string) (pErr error) {
				z, pErr := pipeline.HydrateFromFile(u, a, true)

				if pErr != nil {
					return
				}

				t.Add.PrintZettel(i, z, pErr)

				return
			},
			func(i int, s string, eErr error) {
				t.Add.PrintZettel(i, nil, eErr)
			},
		)

		return
	}
}
