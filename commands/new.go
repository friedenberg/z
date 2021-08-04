package commands

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/friedenberg/z/lib"
	"github.com/friedenberg/z/util"
)

func GetSubcommandNew(f *flag.FlagSet) CommandRunFunc {
	var shouldEdit, shouldPrintFilename, shouldOpenFile bool
	var metadata_json, content, urlToAdd, filePathToAdd string

	f.BoolVar(&shouldEdit, "edit", false, "open the newly created zettel")
	f.BoolVar(&shouldOpenFile, "open-file", false, "open the passed-in file")
	f.BoolVar(&shouldPrintFilename, "print-filename", false, "print the resulting zettel's filename")
	f.StringVar(&urlToAdd, "with-url", "", "include the passed-in URL in the zettel")
	f.StringVar(&filePathToAdd, "with-file", "", "move the passed-in file into zettel control")
	f.StringVar(&content, "content", "", "use the passed-in string as the body. Pass in '-' to read from stdin.")
	f.StringVar(&metadata_json, "metadata-json", "", "parse the passed-in string as the metadata.")

	return func(e *lib.Kasten) (err error) {
		currentTime := time.Now()

		z := &lib.Zettel{Kasten: e}
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

		z.IndexData.Tags = []string{"t-added"}

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

			if shouldOpenFile {
				err = z.Open()
			}

			if err != nil {
				return
			}
		}

		if metadata_json != "" {
			err = json.Unmarshal([]byte(metadata_json), &z.IndexData)

			if err != nil {
				err = fmt.Errorf("parsing metadata json: %w", err)
				return
			}
		}

		if content == "-" {
			var b []byte
			b, err = ioutil.ReadAll(os.Stdin)

			if err != nil {
				return
			}

			z.Data.Body = "\n" + string(b)
		} else {
			z.Data.Body = content
		}

		err = z.Write(func(z *lib.Zettel, errIn error) (errOut error) {
			if errIn != nil {
				if z.HasFile() {
					errOut = os.Remove(z.FilePath())
				}
			}

			return
		})

		if err != nil {
			return
		}

		if shouldEdit {
			err = z.Open()

			if err != nil {
				return
			}
		}

		if shouldPrintFilename {
			fmt.Print(z.Path)
		}

		return z.Edit()
	}
}
