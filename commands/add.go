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

func GetSubcommandAdd(f *flag.FlagSet) CommandRunFunc {
	var metadata_json, kind string
	editActions := printer.Actions(printer.ActionEdit)

	f.Var(&editActions, "actions", "action to perform for the matched zettels")

	f.StringVar(&metadata_json, "metadata-json", "", "parse the passed-in string as the metadata.")
	f.StringVar(&kind, "kind", "", "treat the positional arguments as this kind.")

	return func(e *lib.Kasten) (err error) {
		var add func(z *lib.Zettel, t time.Time, p string) lib.OnZettelWriteFunc

		switch kind {
		case "files":
			add = func(z *lib.Zettel, t time.Time, p string) lib.OnZettelWriteFunc {
				z.IndexData.File = strconv.FormatInt(z.Id, 10) + path.Ext(p)
				return lib.AddFileOnWrite(p)
			}
		case "urls":
			add = func(z *lib.Zettel, t time.Time, p string) lib.OnZettelWriteFunc {
				//TODO normalize
				z.IndexData.Url = p
				return lib.AddUrlOnWrite(p, t)
			}
		default:
			err = fmt.Errorf("unsupported kind: '%s'", kind)
			return
		}

		currentTime := time.Now()

		processor := MakeProcessor(
			e,
			f.Args(),
			&printer.MultiplexingZettelPrinter{
				Printer: &printer.ActionZettelPrinter{
					Kasten:  e,
					Actions: editActions,
				},
			},
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

			onWrite := add(z, currentTime, p)

			err = z.Write(onWrite)

			if err != nil {
				err = fmt.Errorf("failed to write: %w", err)
			}

			return
		}

		err = processor.Run()

		return
	}
}
