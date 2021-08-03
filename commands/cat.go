package commands

import (
	"flag"
	"fmt"
	"sort"

	"github.com/friedenberg/z/commands/printer"
	"github.com/friedenberg/z/lib"
)

type outputFormatFunc func() (printer.ZettelPrinter, ActionFunc)
type outputFormatPrinter struct {
	printer printer.ZettelPrinter
}

var (
	outputFormats    map[string]outputFormatPrinter
	outputFormatKeys []string
)

func init() {
	outputFormats = map[string]outputFormatPrinter{
		"alfred-json": outputFormatPrinter{
			printer: &printer.AlfredJsonZettelPrinter{},
		},
		"alfred-snippet-json": outputFormatPrinter{
			//TODO
			printer: &printer.AlfredJsonZettelPrinter{},
		},
		"metadata-json": outputFormatPrinter{
			printer: &printer.JsonZettelPrinter{},
		},
		"alfred-tags": outputFormatPrinter{
			printer: &printer.Tags{},
		},
		"alfred-expanded-tags": outputFormatPrinter{
			printer: &printer.Tags{ShouldExpand: true},
		},
		"full": outputFormatPrinter{
			printer: &printer.FullZettelPrinter{},
		},
		"filename": outputFormatPrinter{
			printer: &printer.FilenameZettelPrinter{},
		},
	}

	for k, _ := range outputFormats {
		outputFormatKeys = append(outputFormatKeys, k)
	}

	sort.Slice(outputFormatKeys, func(i, j int) bool { return outputFormatKeys[i] < outputFormatKeys[j] })
}

func GetSubcommandCat(f *flag.FlagSet) CommandRunFunc {
	var outputFormat string
	f.StringVar(&outputFormat, "output-format", "full", fmt.Sprintf("One of %q", outputFormatKeys))
	// f.StringVar(&query, "query", "t:snippet", "zettel-spec")

	return func(e *lib.Env) (err error) {
		var p printer.ZettelPrinter
		var actioner ActionFunc

		if format, ok := outputFormats[outputFormat]; ok {
			p = format.printer
		} else {
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
