package commands

import (
	"flag"
	"fmt"
	"path"
	"strconv"
	"time"

	"github.com/friedenberg/z/lib"
)

func GetSubcommandAddFiles(f *flag.FlagSet) CommandRunFunc {
	var shouldEdit, shouldOpen bool

	f.BoolVar(&shouldEdit, "edit", true, "open the created zettel")
	f.BoolVar(&shouldOpen, "open", true, "open the attached file(s)")

	return func(e *lib.Env) (err error) {
		currentTime := time.Now()

		processor := MakeProcessor(
			e,
			f.Args(),
			&nullZettelPrinter{},
		)

		processor.argNormalizer = func(i int, arg string) (normalizedArg string, err error) {
			normalizedArg = arg
			return
		}

		processor.hydrator = func(i int, z *lib.Zettel, p string) (err error) {
			d, err := time.ParseDuration(strconv.Itoa(i) + "s")

			if err != nil {
				panic(err)
			}

			t := currentTime.Add(d)
			z.InitFromTime(t)

			z.IndexData.Tags = []string{"t-added"}
			z.IndexData.File = strconv.FormatInt(z.Id, 10) + path.Ext(p)

			err = z.Write(lib.AddFileOnWrite(p))

			if err != nil {
				err = fmt.Errorf("failed to write: %w", err)
			}

			return
		}

		if shouldEdit {
			processor.actioner = func(i int, z *lib.Zettel) (shouldPrint bool, actionErr error) {
				shouldPrint = true

				if shouldEdit {
					actionErr = z.Edit()
				}

				if err != nil {
					return
				}

				if shouldOpen {
					actionErr = z.Open()
				}

				return
			}
		}

		err = processor.Run()

		return
	}
}
