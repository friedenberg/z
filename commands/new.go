package commands

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/friedenberg/z/commands/printer"
	"github.com/friedenberg/z/lib"
	"github.com/friedenberg/z/util"
)

func GetSubcommandNew(f *flag.FlagSet) CommandRunFunc {
	var metadata_json, content, urlToAdd, filePathToAdd string
	editActions := printer.Actions(printer.ActionEdit)

	f.Var(&editActions, "actions", "action to perform for the matched zettels")
	//TODO convert to action
	f.StringVar(&urlToAdd, "with-url", "", "include the passed-in URL in the zettel")
	f.StringVar(&filePathToAdd, "with-file", "", "move the passed-in file into zettel control")
	f.StringVar(&content, "content", "", "use the passed-in string as the body. Pass in '-' to read from stdin.")
	f.StringVar(&metadata_json, "metadata-json", "", "parse the passed-in string as the metadata.")

	return func(e lib.Umwelt) (err error) {
		currentTime := time.Now()

		z := &lib.Zettel{FilesAndGit: e.FilesAndGit()}
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

		if urlToAdd != "" {
			err = lib.AddUrlOnWrite(urlToAdd, currentTime)(z, nil)

			if err != nil {
				return
			}
		}

		if filePathToAdd != "" {
			err = lib.AddFileOnWrite(filePathToAdd)(z, nil)

			if err != nil {
				return
			}

			if err != nil {
				return
			}
		}

		if metadata_json != "" {
			err = json.Unmarshal([]byte(metadata_json), &z.Metadata)

			if err != nil {
				err = fmt.Errorf("parsing metadata json: %w", err)
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
			Umwelt:  e,
		}

		actionPrinter.Begin()
		actionPrinter.PrintZettel(0, z, nil)
		actionPrinter.End()

		return
	}
}
