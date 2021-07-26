package commands

import (
	"flag"

	"github.com/friedenberg/z/lib"
)

func GetSubcommandCat(f *flag.FlagSet) CommandRunFunc {
	var outputFormat string
	f.StringVar(&outputFormat, "output-format", "full", "One of 'alfred-snippet-json', 'alfred-json', 'metadata-json', 'full', 'filename'")
	// f.StringVar(&query, "query", "t:snippet", "zettel-spec")

	return func(e *lib.Env) (err error) {
		var printer zettelPrinter
		var actioner ActionFunc

		switch outputFormat {
		case "alfred-json":
			printer = &alfredJsonZettelPrinter{}
			format := lib.GetAlfredFormatDefault()
			actioner = func(i int, z *lib.Zettel) error {
				return z.GenerateAlfredItemData(format)
			}

		case "alfred-snippet-json":
			printer = &alfredJsonZettelPrinter{}
			format := lib.GetAlfredFormatSnippet()
			actioner = func(i int, z *lib.Zettel) error {
				return z.GenerateAlfredItemData(format)
			}

		case "metadata-json":
			printer = &jsonZettelPrinter{}
		case "full":
			printer = &fullZettelPrinter{}
		case "filename":
			printer = &filenameZettelPrinter{}
		default:
			printer = &formatZettelPrinter{
				formatter: lib.MakePrintfFormatter(outputFormat),
			}
		}

		processor := MakeProcessor(
			e,
			f.Args(),
			printer,
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
