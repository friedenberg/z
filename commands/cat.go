package commands

import (
	"flag"
	"fmt"

	"github.com/friedenberg/z/lib"
)

func GetSubcommandCat(f *flag.FlagSet) CommandRunFunc {
	var outputFormat string
	f.StringVar(&outputFormat, "output-format", "full", "One of 'alfred-json', 'metadata-json', 'full', 'filename'")

	return func(e *lib.Env) (err error) {
		var printer zettelPrinter
		var actioner ActionFunc

		switch outputFormat {
		case "alfred-json":
			printer = &alfredJsonZettelPrinter{}
			actioner = func(i int, z *lib.Zettel) error {
				return z.GenerateAlfredItemData()
			}

		case "metadata-json":
			printer = &jsonZettelPrinter{}
		case "full":
			printer = &fullZettelPrinter{}
		case "filename":
			printer = &filenameZettelPrinter{}
		default:
			return fmt.Errorf("Unsupported output format: '%s'", outputFormat)
		}

		processor := MakeProcessor(
			e,
			f.Args(),
			printer,
		)

		processor.actioner = actioner

		err = processor.Run()

		return
	}
}
