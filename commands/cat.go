package commands

import (
	"flag"
	"fmt"

	"github.com/friedenberg/z/lib"
	"github.com/friedenberg/z/lib/feeder"
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
	var excludeEmpty bool

	query := filter.MakeAnd()
	tagExclusions := filter.MakeTagExclusions()

	//TODO-P3 rename to "format"
	f.Var(&format, "format", fmt.Sprintf("One of %q", pipeline.FormatKeys))
	f.Var(query, "query", "zettel-spec")
	f.Var(&tagExclusions, "disable-tag-exclusions", "show all zettels, including those excluded by config")
	f.BoolVar(&excludeEmpty, "exclude-empty", true, "don't output empty lines")

	return func(u *lib.Umwelt) (err error) {
		u.ShouldSkipCommit = true

		var args feeder.Feeder

		if len(f.Args()) == 0 {
			args = u.GetAll()
		} else {
			args = feeder.MakeStringSlice(f.Args())
		}

		if excludeEmpty {
			format.SetExcludeEmpty()
		}

		p := pipeline.Pipeline{
			Feeder: args,
			Filter: tagExclusions.WithFilter(
				filter.MakeAnd(query, format.Filter),
				u.TagsForExcludedZettels,
			),
			Writer: format.Writer,
		}

		p.Run(u)

		return
	}
}
