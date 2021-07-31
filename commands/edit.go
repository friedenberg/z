package commands

import (
	"flag"

	"github.com/friedenberg/z/lib"
)

func GetSubcommandEdit(f *flag.FlagSet) CommandRunFunc {
	var shouldEdit bool
	var shouldOpen bool
	var query string

	f.BoolVar(&shouldEdit, "edit", true, "")
	f.BoolVar(&shouldOpen, "open", false, "")
	f.StringVar(&query, "query", "", "zettel-spec string to determine which zettels to open or edit")

	return func(e *lib.Env) (err error) {
		processor := MakeProcessor(
			e,
			f.Args(),
			&multiplexingZettelPrinter{
				printer: &actionZettelPrinter{
					env:        e,
					shouldEdit: shouldEdit,
					shouldOpen: shouldOpen,
				},
			},
		)

		processor.actioner = func(i int, z *lib.Zettel) (shouldPrint bool, err error) {
			shouldPrint = doesZettelMatchQuery(z, query)
			return
		}

		err = processor.Run()

		return
	}
}

func doesZettelMatchQuery(z *lib.Zettel, q string) bool {
	if q == "" {
		return true
	}

	for _, t := range z.IndexData.Tags {
		if t == q {
			return true
		}
	}

	return false
}
