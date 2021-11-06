package commands

import (
	"flag"
	"fmt"

	"github.com/friedenberg/z/lib"
	"github.com/friedenberg/z/lib/pipeline"
	"github.com/friedenberg/z/lib/zettel/filter"
)

func init() {
	makeAndRegisterCommand(
		"cat",
		GetSubcommandCat,
	)
}

func GetSubcommandCat(f *flag.FlagSet) lib.Transactor {
	var format pipeline.Format
	var query string
	//TODO-P3 rename to "format"
	f.Var(&format, "output-format", fmt.Sprintf("One of %q", pipeline.FormatKeys))
	f.StringVar(&query, "query", "", "zettel-spec")
	disableExclude := f.Bool("disable-tag-exclusions", false, "")

	return func(u lib.Umwelt) (err error) {
		u.ShouldSkipCommit = true

		args := f.Args()

		if len(args) == 0 {
			args = u.GetAll()
		}

		p := pipeline.Pipeline{
			Arguments: args,
			Filter: filter.And(
				filter.MatchQuery(query),
				format.Filter,
			),
			Writer: format.Writer,
			//TODO-P3 why was this here?
			// PostWriter: pipeline.MakePostWriter(
			// 	func(i int, z *lib.Zettel) (err error) {
			// 		u.Add.PrintZettel(i, z, err)

			// 		return
			// 	},
			// ),
		}

		if !*disableExclude {
			p.Filter = filter.And(
				p.Filter,
				filter.Not(
					filter.MatchQueries(u.TagsForExcludedZettels...),
				),
			)
		}

		p.Run(u)

		return
	}
}
