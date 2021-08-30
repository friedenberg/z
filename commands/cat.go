package commands

import (
	"flag"
	"fmt"
	"sort"
	"strings"

	"github.com/friedenberg/z/commands/printer"
	"github.com/friedenberg/z/lib"
)

type outputFormatFunc func() (printer.ZettelPrinter, ActionFunc)
type outputFormatPrinter struct {
	printer printer.ZettelPrinter
	filter  func(int, *lib.Zettel) (bool, error)
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
		"alfred-json-files": outputFormatPrinter{
			printer: &printer.AlfredJsonZettelPrinter{
				ItemFunc: printer.AlfredItemsFromZettelFiles,
			},
			filter: func(i int, z *lib.Zettel) (shouldPrint bool, err error) {
				shouldPrint = z.HasFile()
				return
			},
		},
		"alfred-json-urls": outputFormatPrinter{
			printer: &printer.AlfredJsonZettelPrinter{
				ItemFunc: printer.AlfredItemsFromZettelUrls,
			},
			filter: func(i int, z *lib.Zettel) (shouldPrint bool, err error) {
				shouldPrint = z.HasUrl()
				return
			},
		},
		"alfred-json-all": outputFormatPrinter{
			printer: &printer.AlfredJsonZettelPrinter{
				ItemFunc: printer.AlfredItemsFromZettelAll,
			},
		},
		"alfred-json-snippets": outputFormatPrinter{
			//TODO
			printer: &printer.AlfredJsonZettelPrinter{
				ItemFunc: printer.AlfredItemsFromZettelSnippets,
			},
			filter: func(i int, z *lib.Zettel) (shouldPrint bool, err error) {
				for _, t := range z.IndexData.Tags {
					if strings.Contains(t, "t-snippet") {
						shouldPrint = true
					}
				}

				return
			},
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
	var outputFormat, query string
	f.StringVar(&outputFormat, "output-format", "full", fmt.Sprintf("One of %q", outputFormatKeys))
	f.StringVar(&query, "query", "", "zettel-spec")

	return func(e *lib.Kasten) (err error) {
		var p printer.ZettelPrinter
		var actioner ActionFunc

		if format, ok := outputFormats[outputFormat]; ok {
			p = format.printer
			actioner = format.filter
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

		processor.actioner = func(i int, z *lib.Zettel) (shouldPrint bool, err error) {
			if actioner != nil {
				shouldPrint, err = actioner(i, z)
				shouldPrint = shouldPrint && doesZettelMatchQuery(z, query)
			} else {
				shouldPrint = doesZettelMatchQuery(z, query)
			}

			return
		}

		err = processor.Run()

		return
	}
}
