package commands

import (
	"bufio"
	"flag"
	"fmt"
	"os"

	"github.com/friedenberg/z/lib"
	"github.com/friedenberg/z/lib/pipeline"
	"github.com/friedenberg/z/lib/zettel/modifier"
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
		var args []string

		if stdin {
			args = make([]string, 0)
			scanner := bufio.NewScanner(os.Stdin)

			for scanner.Scan() {
				args = append(args, scanner.Text())
			}
		} else {
			args = f.Args()
		}

		p := pipeline.Pipeline{
			Arguments: args,
			Reader:    format.Reader,
			Modifier:  modifier.TransactionAction(u.Transaction, lib.TransactionActionAdded),
		}

		p.Run(u)

		return
	}
}
