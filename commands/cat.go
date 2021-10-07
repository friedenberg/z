package commands

import (
	"flag"
	"fmt"
	"sort"
	"strings"

	"github.com/friedenberg/z/lib"
	"github.com/friedenberg/z/lib/pipeline"
	"github.com/friedenberg/z/lib/pipeline/printer"
	"github.com/friedenberg/z/util"
)

type outputFormat pipeline.FilterPrinter

var (
	outputFormats    map[string]outputFormat
	outputFormatKeys []string
)

func init() {
	outputFormats = map[string]outputFormat{
		"alfred-json": outputFormat{
			Printer: &printer.AlfredJsonZettelPrinter{},
		},
		"alfred-json-files": outputFormat{
			Filter: func(i int, z *lib.Zettel) bool {
				return z.HasFile()
			},
			Printer: &printer.AlfredJsonZettelPrinter{
				ItemFunc: printer.AlfredItemsFromZettelFiles,
			},
		},
		"alfred-json-urls": outputFormat{
			Filter: func(i int, z *lib.Zettel) bool {
				return z.HasUrl()
			},
			Printer: &printer.AlfredJsonZettelPrinter{
				ItemFunc: printer.AlfredItemsFromZettelUrls,
			},
		},
		"alfred-json-all": outputFormat{
			Printer: &printer.AlfredJsonZettelPrinter{
				ItemFunc: printer.AlfredItemsFromZettelAll,
			},
		},
		"alfred-json-snippets": outputFormat{
			Filter: func(i int, z *lib.Zettel) bool {
				for _, t := range z.Metadata.Tags {
					if strings.Contains(t, "t-snippet") {
						return true
					}
				}

				return false
			},
			//TODO
			Printer: &printer.AlfredJsonZettelPrinter{
				ItemFunc: printer.AlfredItemsFromZettelSnippets,
			},
		},
		"metadata-json": outputFormat{
			Printer: &printer.JsonZettelPrinter{},
		},
		"alfred-tags": outputFormat{
			Printer: &printer.Tags{},
		},
		"alfred-expanded-tags": outputFormat{
			Printer: &printer.Tags{ShouldExpand: true},
		},
		"full": outputFormat{
			Printer: &printer.FullZettelPrinter{},
		},
		"filename": outputFormat{
			Printer: &printer.FilenameZettelPrinter{},
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
				Printer: &printer.FullZettelPrinter{},
			}
		} else {
			*a = outputFormat{
				Printer: &printer.FormatZettelPrinter{
					Formatter: lib.MakePrintfFormatter(s),
				},
			}
		}
	}

	a.Printer = &printer.MultiplexingZettelPrinter{Printer: a.Printer}

	return
}

func GetSubcommandCat(f *flag.FlagSet) lib.Transactor {
	var of outputFormat
	var query string
	f.Var(&of, "output-format", fmt.Sprintf("One of %q", outputFormatKeys))
	f.StringVar(&query, "query", "", "zettel-spec")

	return func(u lib.Umwelt, t lib.Transaction) (err error) {
		args := f.Args()
		var iter util.ParallelizerIterFunc

		if u.Config.UseIndexCache {
			if len(args) == 0 {
				args = u.GetAll()
			}

			iter = cachedIteration(u, query, pipeline.FilterPrinter(of))
		} else {
			if len(args) == 0 {
				args, err = u.FilesAndGit().GetAll()

				if err != nil {
					return
				}
			}

			iter = filesystemIteration(u, query, pipeline.FilterPrinter(of))
		}

		par := util.Parallelizer{Args: args}
		of.Printer.Begin()
		defer of.Printer.End()
		par.Run(iter, errIterartion(of.Printer))

		return
	}
}
