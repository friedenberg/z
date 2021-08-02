package commands

import (
	"flag"

	"github.com/friedenberg/z/commands/printer"
	"github.com/friedenberg/z/lib"
)

func GetSubcommandCat(f *flag.FlagSet) CommandRunFunc {
	var outputFormat string
	f.StringVar(&outputFormat, "output-format", "full", "One of 'alfred-snippet-json', 'alfred-json', 'metadata-json', 'full', 'filename'")
	// f.StringVar(&query, "query", "t:snippet", "zettel-spec")

	return func(e *lib.Env) (err error) {
		var p printer.ZettelPrinter
		var actioner ActionFunc

		switch outputFormat {
		case "alfred-json":
			p = &printer.AlfredJsonZettelPrinter{}

		case "alfred-snippet-json":
			//TODO add formatter to printer
			p = &printer.AlfredJsonZettelPrinter{}

		case "metadata-json":
			p = &printer.JsonZettelPrinter{}
		case "full":
			p = &printer.FullZettelPrinter{}
		case "filename":
			p = &printer.FilenameZettelPrinter{}
		default:
			p = &printer.FormatZettelPrinter{
				Formatter: lib.MakePrintfFormatter(outputFormat),
			}
		}

		processor := MakeProcessor(
			e,
			f.Args(),
			p,
		)

		processor.hydrator = func(_ int, z *lib.Zettel, path string) error {
			z.Path = path
			return z.HydrateFromFilePath(true)
		}

		processor.actioner = actioner

		err = processor.Run()

		return
	}
}
