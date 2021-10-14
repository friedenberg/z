package commands

import (
	"flag"
	"fmt"

	"github.com/friedenberg/z/lib"
	"github.com/friedenberg/z/lib/pipeline"
	"github.com/friedenberg/z/util"
)

func init() {
	n := "cat"
	f := flag.NewFlagSet(n, flag.ExitOnError)

	registerCommand(
		n,
		Command{
			Flags: f,
			Run:   GetSubcommandCat(f),
		},
	)
}

func GetSubcommandCat(f *flag.FlagSet) lib.Transactor {
	var of outputFormat
	var query string
	f.Var(&of, "output-format", fmt.Sprintf("One of %q", outputFormatKeys))
	f.StringVar(&query, "query", "", "zettel-spec")

	return func(u lib.Umwelt) (err error) {
		u.ShouldSkipCommit = true
		args := f.Args()
		var iter util.ParallelizerIterFunc

		if len(args) == 0 {
			args = u.GetAll()
		}

		iter = cachedIteration(u, query, pipeline.FilterPrinter(of))

		par := util.Parallelizer{Args: args}
		of.Printer.Begin()
		defer of.Printer.End()
		par.Run(iter, errIterartion(of.Printer))

		return
	}
}
