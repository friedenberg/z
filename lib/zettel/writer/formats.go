package writer

import (
	"sort"

	"github.com/friedenberg/z/lib/pipeline"
	"github.com/friedenberg/z/lib/zettel/filter"
)

type Format pipeline.Pipeline

var (
	Formats    map[string]Format
	FormatKeys []string
)

func init() {
	Formats = map[string]Format{
		"alfred-json": Format{
			Writer: &AlfredJson{},
		},
		"alfred-json-files": Format{
			Filter: filter.HasFile(),
			Writer: &AlfredJson{
				ItemFunc: AlfredItemsFromZettelFiles,
			},
		},
		"alfred-json-urls": Format{
			Filter: filter.HasUrl(),
			Writer: &AlfredJson{
				ItemFunc: AlfredItemsFromZettelUrls,
			},
		},
		"alfred-json-all": Format{
			Writer: &AlfredJson{
				ItemFunc: AlfredItemsFromZettelAll,
			},
		},
		"alfred-json-snippets": Format{
			Filter: filter.MatchQuery("t-snippet"),
			//TODO
			Writer: &AlfredJson{
				ItemFunc: AlfredItemsFromZettelSnippets,
			},
		},
		"alfred-tags": Format{
			Writer: &Tags{},
		},
		"alfred-expanded-tags": Format{
			Writer: &Tags{ShouldExpand: true},
		},
		"full": Format{
			Writer: &Full{},
		},
		"filename": Format{
			Writer: MakeFormatter("%p"),
		},
		"toml-to-json": Format{
			// Filter: filter.MatchQuery("k-toml"),
			Writer: &TomlToJson{},
		},
		"try-format": Format{
			Filter: filter.Or(
				filter.MatchQuery("k-toml"),
				filter.MatchQuery("from-pb"),
			),
			Writer: &TryFormat{},
		},
		"json": Format{
			Reader: &Json{},
			Writer: &Json{},
		},
		"json-with-body": Format{
			Writer: &Json{IncludeBody: true},
		},
	}

	for k, _ := range Formats {
		FormatKeys = append(FormatKeys, k)
	}

	sort.Slice(FormatKeys, func(i, j int) bool { return FormatKeys[i] < FormatKeys[j] })
}

func (a *Format) String() string {
	//TODO
	return ""
}

func (a *Format) Set(s string) (err error) {
	if format, ok := Formats[s]; ok {
		*a = format
	} else {
		if s == "" {
			*a = Format{
				Writer: &Full{},
			}
		} else {
			*a = Format{
				Writer: MakeFormatter(s),
			}
		}
	}

	return
}
