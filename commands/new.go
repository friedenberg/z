package commands

import (
	"flag"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/friedenberg/z/lib"
)

func GetSubcommandNew(f *flag.FlagSet) CommandRunFunc {
	var shouldOpen, shouldPrintFilename bool
	var urlToAdd, filePathToAdd string

	f.BoolVar(&shouldOpen, "open", false, "open the newly created zettel")
	f.BoolVar(&shouldPrintFilename, "print-filename", false, "print the resulting zettel's filename")
	f.StringVar(&urlToAdd, "with-url", "", "include the passed-in URL in the zettel")
	f.StringVar(&filePathToAdd, "with-file", "", "move the passed-in file into zettel control")

	return func(e Env) (err error) {
		currentTime := time.Now()

		z := &lib.Zettel{}
		z.InitFromTime(e.ZettelPath, currentTime)

		z.IndexData.Tags = []string{"added"}

		zi := path.Base(z.Path)[:len(path.Ext(z.Path))]

		if urlToAdd != "" {
			err = lib.AddUrlOnWrite(urlToAdd, currentTime)(z, nil)

			if err != nil {
				return
			}
		}

		if filePathToAdd != "" {
			err = lib.AddFileOnWrite(e.ZettelPath, filePathToAdd, zi)(z, nil)

			if err != nil {
				return
			}
		}

		err = z.Write(func(z *lib.Zettel, errIn error) (errOut error) {
			if errIn != nil {
				if z.HasFile() {
					errOut = os.Remove(z.IndexData.File)
				}
			}

			return
		})

		if err != nil {
			return
		}

		if shouldOpen {
			z.Open(e.ZettelPath)
		}

		if shouldPrintFilename {
			fmt.Print(z.Path)
		}

		return
	}
}
