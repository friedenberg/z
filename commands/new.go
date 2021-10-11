package commands

import (
	"flag"
	"io/ioutil"
	"os"
	"strings"

	"github.com/friedenberg/z/commands/options"
	"github.com/friedenberg/z/lib"
	"github.com/friedenberg/z/lib/pipeline"
	"github.com/friedenberg/z/lib/pipeline/printer"
)

func GetSubcommandNew(f *flag.FlagSet) lib.Transactor {
	var tags, content string
	editActions := options.Actions(options.ActionEdit)

	f.Var(&editActions, "actions", "action to perform for the matched zettels")
	f.StringVar(&content, "content", "", "use the passed-in string as the body. Pass in '-' to read from stdin.")
	f.StringVar(&tags, "tags", "", "use the passed-in space-separated string as tags")

	return func(u lib.Umwelt, t *lib.Transaction) (err error) {
		z, err := pipeline.New(u)

		if err != nil {
			return
		}

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

		//TODO-P3 add to transaction
		err = z.Write(func(z *lib.Zettel, errIn error) (errOut error) {
			if errIn != nil {
				if f, ok := z.Note.Metadata.LocalFile(); ok {
					errOut = os.Remove(f.FilePath(u.BasePath))
				}

				return
			}

			return
		})

		if err != nil {
			return
		}

		actionPrinter := printer.ActionZettelPrinter{
			Actions:     editActions,
			Umwelt:      u,
			Transaction: t,
		}

		actionPrinter.Begin()
		actionPrinter.PrintZettel(0, z, nil)
		actionPrinter.End()

		t.Add.PrintZettel(0, z, nil)

		return
	}
}
