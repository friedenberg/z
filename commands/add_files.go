package commands

import (
	"encoding/json"
	"flag"
	"fmt"
	"path"
	"strconv"
	"time"

	"github.com/friedenberg/z/commands/printer"
	"github.com/friedenberg/z/lib"
	"github.com/friedenberg/z/util"
)

func GetSubcommandAddFiles(f *flag.FlagSet) CommandRunFunc {
	var shouldEdit, shouldOpen bool
	var metadata_json string

	f.BoolVar(&shouldEdit, "edit", true, "open the created zettel")
	f.BoolVar(&shouldOpen, "open", true, "open the attached file(s)")

	f.StringVar(&metadata_json, "metadata-json", "", "parse the passed-in string as the metadata.")

	return func(e *lib.Kasten) (err error) {
		currentTime := time.Now()

		processor := MakeProcessor(
			e,
			f.Args(),
			&printer.NullZettelPrinter{},
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

			if metadata_json != "" {
				err = json.Unmarshal([]byte(metadata_json), &z.IndexData)

				if err != nil {
					err = fmt.Errorf("parsing metadata json: %w", err)
					return
				}
			}

			z.IndexData.Tags = append(z.IndexData.Tags, "zz-inbox")

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
