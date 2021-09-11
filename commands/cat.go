package commands

import (
	"flag"
	"fmt"
	"sort"
	"strings"

	"github.com/friedenberg/z/commands/printer"
	"github.com/friedenberg/z/lib"
	"github.com/friedenberg/z/lib/pipeline"
	"github.com/friedenberg/z/util"
)

type outputFormat struct {
	printer printer.ZettelPrinter
	filter  func(int, *lib.Zettel) bool
}

var (
	outputFormats    map[string]outputFormat
	outputFormatKeys []string
)

func init() {
	outputFormats = map[string]outputFormat{
		"alfred-json": outputFormat{
			printer: &printer.AlfredJsonZettelPrinter{},
		},
		"alfred-json-files": outputFormat{
			printer: &printer.AlfredJsonZettelPrinter{
				ItemFunc: printer.AlfredItemsFromZettelFiles,
			},
			filter: func(i int, z *lib.Zettel) bool {
				return z.HasFile()
			},
		},
		"alfred-json-urls": outputFormat{
			printer: &printer.AlfredJsonZettelPrinter{
				ItemFunc: printer.AlfredItemsFromZettelUrls,
			},
			filter: func(i int, z *lib.Zettel) bool {
				return z.HasUrl()
			},
		},
		"alfred-json-all": outputFormat{
			printer: &printer.AlfredJsonZettelPrinter{
				ItemFunc: printer.AlfredItemsFromZettelAll,
			},
		},
		"alfred-json-snippets": outputFormat{
			//TODO
			printer: &printer.AlfredJsonZettelPrinter{
				ItemFunc: printer.AlfredItemsFromZettelSnippets,
			},
			filter: func(i int, z *lib.Zettel) bool {
				for _, t := range z.Metadata.Tags {
					if strings.Contains(t, "t-snippet") {
						return true
					}
				}

				return false
			},
		},
		"metadata-json": outputFormat{
			printer: &printer.JsonZettelPrinter{},
		},
		"alfred-tags": outputFormat{
			printer: &printer.Tags{},
		},
		"alfred-expanded-tags": outputFormat{
			printer: &printer.Tags{ShouldExpand: true},
		},
		"full": outputFormat{
			printer: &printer.FullZettelPrinter{},
		},
		"filename": outputFormat{
			printer: &printer.FilenameZettelPrinter{},
		},
	}

	for k, _ := range outputFormats {
		outputFormatKeys = append(outputFormatKeys, k)
	}

	sort.Slice(outputFormatKeys, func(i, j int) bool { return outputFormatKeys[i] < outputFormatKeys[j] })
}

func (a *outputFormat) String() string {
	//TODO
	return ""
}

func (a *outputFormat) Set(s string) (err error) {
	if format, ok := outputFormats[s]; ok {
		*a = format
	} else {
		if s == "" {
			*a = outputFormat{
				printer: &printer.FullZettelPrinter{},
			}
		} else {
			*a = outputFormat{
				printer: &printer.FormatZettelPrinter{
					Formatter: lib.MakePrintfFormatter(s),
				},
			}
		}
	}

	a.printer = &printer.MultiplexingZettelPrinter{Printer: a.printer}

	return
}

func GetSubcommandCat(f *flag.FlagSet) CommandRunFunc {
	var of outputFormat
	var query string
	f.Var(&of, "output-format", fmt.Sprintf("One of %q", outputFormatKeys))
	f.StringVar(&query, "query", "", "zettel-spec")

	return func(e lib.Umwelt) (err error) {
		args := f.Args()
		var iter util.ParallelizerIterFunc

		if e.Config.UseIndexCache {
			if len(args) == 0 {
				args = e.GetAll()
			}

			iter = cachedIteration(e, query, of)
		} else {
			if len(args) == 0 {
				args, err = e.FilesAndGit().GetAll()

				if err != nil {
					return
				}
			}

			iter = filesystemIteration(e, query, of)
		}

		par := util.Parallelizer{Args: args}
		of.printer.Begin()
		defer of.printer.End()
		par.Run(iter, errIterartion(of))

		return
	}
}

func printIfNecessary(u lib.Umwelt, i int, z *lib.Zettel, q string, o outputFormat) {
	if (o.filter == nil || o.filter(i, z)) && doesZettelMatchQuery(z, q) {
		o.printer.PrintZettel(i, z, nil)
	}
}

func cachedIteration(u lib.Umwelt, q string, o outputFormat) util.ParallelizerIterFunc {
	return func(i int, s string) (err error) {
		z, err := pipeline.HydrateFromIndex(u, s)

		if err != nil {
			return
		}

		printIfNecessary(u, i, z, q, o)

		return
	}
}

func filesystemIteration(u lib.Umwelt, q string, o outputFormat) util.ParallelizerIterFunc {
	return func(i int, s string) (err error) {
		p, err := pipeline.NormalizePath(u, s)

		if err != nil {
			return
		}
		//TODO determine if body read is necessary
		z, err := pipeline.HydrateFromFile(u, p, true)

		if err != nil {
			return
		}

		printIfNecessary(u, i, z, q, o)

		return
	}
}

func errIterartion(o outputFormat) util.ParallelizerErrorFunc {
	return func(i int, s string, err error) {
		o.printer.PrintZettel(i, nil, err)
	}
}
