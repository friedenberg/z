package commands

import (
	"flag"

	"github.com/friedenberg/z/lib"
	"github.com/friedenberg/z/lib/pipeline"
	"github.com/friedenberg/z/util"
)

func init() {
	makeAndRegisterCommand(
		"index",
		GetSubcommandIndex,
	)
}

func GetSubcommandIndex(f *flag.FlagSet) lib.Transactor {
	return func(u lib.Umwelt) (err error) {
		u.ShouldSkipCommit = true
		u.Index = lib.MakeIndex()

		args, err := u.Store().GetAll()

		if err != nil {
			return
		}

		par := util.Parallelizer{Args: args}
		u.Add.Begin()
		defer u.Add.End()
		par.Run(
			func(i int, a string) (pErr error) {
				z, pErr := pipeline.HydrateFromFile(u, a, true)

				if pErr != nil {
					return
				}

				u.Add.PrintZettel(i, z, pErr)

				return
			},
			func(i int, s string, eErr error) {
				u.Add.PrintZettel(i, nil, eErr)
			},
		)

		return
	}
}
