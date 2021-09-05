package commands

import (
	"flag"
	"fmt"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/friedenberg/z/commands/printer"
	"github.com/friedenberg/z/lib"
)

func GetSubcommandAdd(f *flag.FlagSet) CommandRunFunc {
	var tagString, kind string
	editActions := printer.Actions(printer.ActionEdit)

	f.Var(&editActions, "actions", "action to perform for the matched zettels")

	f.StringVar(&tagString, "tags", "", "parse the passed-in string as the metadata.")
	f.StringVar(&kind, "kind", "", "treat the positional arguments as this kind.")

	return func(e *lib.Kasten) (err error) {
		currentTime := time.Now()

		bootstrapZettel := func(i int, z *lib.Zettel, p string) (err error) {
			if z.Id == 0 {
				err = z.InitAndAssignUniqueId(currentTime, i)

				if err != nil {
					panic(err)
				}
			}

			return
		}

		var add func(z *lib.Zettel, t time.Time, p string) lib.OnZettelWriteFunc
		var hydrator func(i int, z *lib.Zettel, p string) (err error)

		switch kind {
		case "files":
			add = func(z *lib.Zettel, t time.Time, p string) lib.OnZettelWriteFunc {
				z.Metadata.File = strconv.FormatInt(z.Id, 10) + path.Ext(p)
				return lib.AddFileOnWrite(p)
			}
			hydrator = bootstrapZettel
		case "urls":
			add = func(z *lib.Zettel, t time.Time, p string) lib.OnZettelWriteFunc {
				//TODO normalize
				z.Metadata.Url = p
				return lib.AddUrlOnWrite(p, t)
			}
			hydrator = func(i int, z *lib.Zettel, p string) (err error) {
				indexItems := e.Index.ZettelsForUrl(p)

				if len(indexItems) > 1 {
					err = fmt.Errorf("multiple zettels ('%q') with url: '%s'", indexItems, p)
					return
				}

				if len(indexItems) == 1 {
					e.Index.HydrateZettel(z, indexItems[0])
				}

				err = bootstrapZettel(i, z, p)
				return
			}
		default:
			err = fmt.Errorf("unsupported kind: '%s'", kind)
			return
		}

		err = hydrateIndex(e)

		if err != nil {
			return
		}

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
			err = hydrator(i, z, p)

			if err != nil {
				return
			}

			if tagString != "" {
				tags := strings.Split(tagString, " ")
				z.Metadata.Tags = append(z.Metadata.Tags, tags...)
			} else {
				z.Metadata.Tags = append(z.Metadata.Tags, "zz-inbox")
			}

			z.Metadata.Tags = uniqueAndSortTags(z.Metadata.Tags)

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
