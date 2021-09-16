package commands

import (
	"flag"

	"github.com/friedenberg/z/lib"
	"github.com/friedenberg/z/lib/pipeline"
	"golang.org/x/xerrors"
)

func init() {
	makeAndRegisterCommand(
		"rm",
		GetSubcommandRm,
	)
}

func GetSubcommandRm(f *flag.FlagSet) lib.Transactor {
	return func(u lib.Umwelt) (err error) {
		args := f.Args()

		if len(args) == 0 {
			err = xerrors.Errorf("no zettels included for deletion")
			return
		}

		p := pipeline.Pipeline{
			Arguments: args,
			Modifier:  u.Del,
		}

		p.Run(u)

		return
	}
}
