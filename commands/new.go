package commands

import (
	"flag"
	"io/ioutil"
	"os"
	"strings"

	"github.com/friedenberg/z/commands/options"
	"github.com/friedenberg/z/lib"
	"github.com/friedenberg/z/lib/pipeline"
	"github.com/friedenberg/z/lib/zettel"
	"github.com/friedenberg/z/lib/zettel/modifier"
	"github.com/friedenberg/z/lib/zettel/reader"
)

func init() {
	makeAndRegisterCommand(
		"new",
		GetSubcommandNew,
	)
}

func GetSubcommandNew(f *flag.FlagSet) lib.Transactor {
	var tags, content string
	editActions := options.Actions(options.ActionEdit)

	f.Var(&editActions, "actions", "action to perform for the matched zettels")
	f.StringVar(&content, "content", "", "use the passed-in string as the body. Pass in '-' to read from stdin.")
	f.StringVar(&tags, "tags", "", "use the passed-in space-separated string as tags")

	return func(u *lib.Umwelt) (err error) {
		p := pipeline.Pipeline{
			Arguments: []string{""},
			Reader:    reader.New(),
			Modifier: modifier.Chain(
				modifier.Make(
					func(i int, z *zettel.Zettel) (err error) {
						if tags != "" {
							z.Note.Metadata.SetStringTags(strings.Split(tags, " "))
						}

						if content == "-" {
							var b []byte
							b, err = ioutil.ReadAll(os.Stdin)

							if err != nil {
								return
							}

							z.Body = "\n" + string(b)
						} else {
							z.Body = content
						}

						return
					},
				),
				modifier.Make(
					func(i int, z *zettel.Zettel) (err error) {
						err = z.Write(nil)
						return
					},
				),
				&lib.ActionModifier{
					Umwelt:  u,
					Actions: editActions,
				},
				lib.MakeTransactionAction(u.Transaction, lib.TransactionActionAdded),
			),
		}

		p.Run(u)

		return
	}
}
