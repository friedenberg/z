package commands

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"os"
	"time"

	"github.com/friedenberg/z/commands/options"
	"github.com/friedenberg/z/commands/printer"
	"github.com/friedenberg/z/lib"
	"github.com/friedenberg/z/util"
	"golang.org/x/xerrors"
)

func GetSubcommandNew(f *flag.FlagSet) CommandRunFunc {
	var metadata_json, content string
	editActions := options.Actions(options.ActionEdit)

	f.Var(&editActions, "actions", "action to perform for the matched zettels")
	f.StringVar(&content, "content", "", "use the passed-in string as the body. Pass in '-' to read from stdin.")
	f.StringVar(&metadata_json, "metadata-json", "", "parse the passed-in string as the metadata.")

	return func(u lib.Umwelt) (err error) {
		currentTime := time.Now()

		z := &lib.Zettel{
			Umwelt: &u,
		}
		z.InitFromTime(currentTime)

		for {
			if util.FileExists(z.Path) {
				d, err := time.ParseDuration("1s")

				if err != nil {
					panic(err)
				}

				currentTime = currentTime.Add(d)
				z.InitFromTime(currentTime)
			} else {
				break
			}
		}

		if err != nil {
			return
		}

		if metadata_json != "" {
			err = json.Unmarshal([]byte(metadata_json), &z.Metadata)

			if err != nil {
				err = xerrors.Errorf("parsing metadata json: %w", err)
				return
			}
		}

		z.Metadata.Tags = append(z.Metadata.Tags, "zz-inbox")

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

		err = z.Write(func(z *lib.Zettel, errIn error) (errOut error) {
			if errIn != nil {
				if z.HasFile() {
					errOut = os.Remove(z.FilePath())
				}

				return
			}

			return
		})

		if err != nil {
			return
		}

		actionPrinter := printer.ActionZettelPrinter{
			Actions: editActions,
			Umwelt:  u,
		}

		actionPrinter.Begin()
		actionPrinter.PrintZettel(0, z, nil)
		actionPrinter.End()

		return
	}
}
