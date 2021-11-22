package commands

import (
	"flag"
	"fmt"
	"os"

	"github.com/friedenberg/z/lib"
	"github.com/friedenberg/z/lib/feeder"
	"github.com/friedenberg/z/lib/pipeline"
)

func init() {
	makeAndRegisterCommand(
		"import",
		GetSubcommandImport,
	)
}

func GetSubcommandImport(f *flag.FlagSet) lib.Transactor {
	var format pipeline.Format
	var stdin bool
	f.Var(&format, "format", fmt.Sprintf("One of %q", pipeline.FormatKeys))
	f.BoolVar(&stdin, "stdin", false, "use stdin for input")

	return func(u *lib.Umwelt) (err error) {
		var args feeder.Feeder

		if stdin {
			args = feeder.MakeIoReader(os.Stdin)
		} else {
			args = feeder.MakeStringSlice(f.Args())
		}

		p := pipeline.Pipeline{
			Feeder:   args,
			Reader:   format.Reader,
			Modifier: lib.MakeTransactionAction(u.Transaction, lib.TransactionActionAdded),
		}

		p.Run(u)

		return
	}
}
