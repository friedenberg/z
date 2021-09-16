package commands

import (
	"flag"
	"strings"

	"github.com/friedenberg/z/commands/options"
	"github.com/friedenberg/z/lib"
	"github.com/friedenberg/z/lib/pipeline"
	"github.com/friedenberg/z/lib/pipeline/modifier"
)

func init() {
	makeAndRegisterCommand(
		"add",
		GetSubcommandAdd,
	)
}

func GetSubcommandAdd(f *flag.FlagSet) lib.Transactor {
	var tagString string
	var description string
	var kind attachmentKind
	editActions := options.Actions(options.ActionEdit)

	f.Var(&editActions, "actions", "action to perform for the matched zettels")

	f.StringVar(&tagString, "tags", "", "parse the passed-in string as the metadata.")
	f.StringVar(&description, "description", "", "use this string as the zettel description")
	f.Var(&kind, "kind", "treat the positional arguments as this kind.")

	return func(u lib.Umwelt) (err error) {
		p := pipeline.Pipeline{
			Arguments: f.Args(),
			Reader:    kind,
			Modifier: modifier.Chain(
				modifier.Make(
					func(i int, z *lib.Zettel) (err error) {
						tags := strings.Split(tagString, " ")
						z.Metadata.AddStringTags(tags...)

						if description != "" {
							z.Metadata.SetDescription(description)
						}

						return
					},
				),
				u.Add,
			),
		}

		p.Run(u)

		added := u.Added().Paths()

		err = u.RunTransaction(nil)

		if err != nil {
			return
		}

		//TODO-P4 check why this is re-using added zettels rather than modifying new
		//ones
		u.Transaction = lib.MakeTransaction()

		p = pipeline.Pipeline{
			Arguments: added,
			Modifier: &modifier.Action{
				Umwelt:  u,
				Actions: editActions,
			},
		}

		p.Run(u)

		return
	}
}
