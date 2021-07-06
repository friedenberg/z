package commands

import (
	"flag"
	"fmt"
	"strconv"
	"time"

	"github.com/friedenberg/z/lib"
)

func GetSubcommandAdd(f *flag.FlagSet) CommandRunFunc {
	isUrl := false
	f.BoolVar(&isUrl, "url", false, "")

	return func(e Env) (err error) {
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
			z.InitFromTime(e.ZettelPath, t)

			//TODO move into zettel init
			zettelId := strconv.FormatInt(t.Unix(), 10)
			z.IndexData.Tags = []string{"added"}

			var onWrite lib.OnZettelWriteFunc

			if isUrl {
				onWrite = lib.AddUrlOnWrite(p, t)
			} else {
				onWrite = lib.AddFileOnWrite(e.ZettelPath, p, zettelId)
			}

			if err != nil {
				err = fmt.Errorf("failed to add url or file: %w", err)
				return
			}

			err = z.Write(onWrite)

			if onWrite != nil {
				err = onWrite(z, err)
			}

			if err != nil {
				err = fmt.Errorf("failed to write: %w", err)
			}

			return
		}

		err = processor.Run()

		return
	}
}
