package pipeline

import (
	"sort"

	"github.com/friedenberg/z/lib/zettel/filter"
	"github.com/friedenberg/z/lib/zettel/writer"
)

type Format Pipeline

var (
	Formats    map[string]Format
	FormatKeys []string
)

func init() {
	Formats = map[string]Format{
		"alfred-json": Format{
			Writer: &writer.AlfredJson{},
		},
		"alfred-json-files": Format{
			Filter: filter.HasFile(),
			Writer: &writer.AlfredJson{
				ItemFunc: writer.AlfredItemsFromZettelFiles,
			},
		},
		"alfred-json-urls": Format{
			Filter: filter.HasUrl(),
			Writer: &writer.AlfredJson{
				ItemFunc: writer.AlfredItemsFromZettelUrls,
			},
		},
		"alfred-json-all": Format{
			Writer: &writer.AlfredJson{
				ItemFunc: writer.AlfredItemsFromZettelAll,
			},
		},
		"alfred-json-snippets": Format{
			Filter: filter.MatchQuery("t-snippet"),
			//TODO
			Writer: &writer.AlfredJson{
				ItemFunc: writer.AlfredItemsFromZettelSnippets,
			},
		},
		"alfred-tags": Format{
			Writer: &writer.Tags{},
		},
		"alfred-expanded-tags": Format{
			Writer: &writer.Tags{ShouldExpand: true},
		},
		"full": Format{
			Writer: &writer.Full{},
		},
		"filename": Format{
			Writer: writer.MakeFormatter("%p"),
		},
		"toml-to-json": Format{
			// Filter: filter.MatchQuery("k-toml"),
			Writer: &writer.TomlToJson{},
		},
		"try-format": Format{
			Filter: filter.Or(
				filter.MatchQuery("k-toml"),
				filter.MatchQuery("from-pb"),
			),
			Writer: &writer.TryFormat{},
		},
		"json": Format{
			Reader: &writer.Json{},
			Writer: &writer.Json{},
		},
		"json-with-body": Format{
			Writer: &writer.Json{IncludeBody: true},
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
				Writer: &writer.Full{},
			}
		} else {
			*a = Format{
				Writer: writer.MakeFormatter(s),
			}
		}
	}

	return
}
