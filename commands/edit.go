package commands

import (
	"flag"
	"strings"

	"github.com/friedenberg/z/commands/printer"
	"github.com/friedenberg/z/lib"
)

func GetSubcommandEdit(f *flag.FlagSet) CommandRunFunc {
	var shouldEdit bool
	var shouldOpen bool
	var query string

	f.BoolVar(&shouldEdit, "edit", true, "")
	f.BoolVar(&shouldOpen, "open", false, "")
	f.StringVar(&query, "query", "", "zettel-spec string to determine which zettels to open or edit")

	return func(e *lib.Kasten) (err error) {
		processor := MakeProcessor(
			e,
			f.Args(),
			&printer.MultiplexingZettelPrinter{
				Printer: &printer.ActionZettelPrinter{
					Kasten:        e,
					ShouldEdit: shouldEdit,
					ShouldOpen: shouldOpen,
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
		if strings.Contains(t, q) {
			return true
		}
	}

	return false
}
