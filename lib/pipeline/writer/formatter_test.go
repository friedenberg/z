package writer

import (
	"testing"

	"github.com/friedenberg/z/lib"
)

type printFormatterTestCaseMakeZettelFunc func() *lib.Zettel

type printfFormatterTestCase struct {
	name       string
	makeZettel printFormatterTestCaseMakeZettelFunc
	format     string
	output     string
}

func getPrintfTestCases(t *testing.T) []printfFormatterTestCase {
	k := lib.FileStore{}
	umwelt := lib.Umwelt{
		Kasten: &k,
	}

	makeZettelWithDate := func() (z *lib.Zettel) {
		z = &lib.Zettel{
			Umwelt: umwelt,
			Path:   "1633902356.md",
		}

		return
	}

	return []printfFormatterTestCase{
		printfFormatterTestCase{
			name:       "just date",
			makeZettel: makeZettelWithDate,
			format:     "%w",
			output:     "2021-10-10",
		},
		printfFormatterTestCase{
			name:       "just constant string",
			makeZettel: makeZettelWithDate,
			format:     "wow",
			output:     "wow",
		},
		printfFormatterTestCase{
			name:       "just constant one-char string",
			makeZettel: makeZettelWithDate,
			format:     "w",
			output:     "w",
		},
		printfFormatterTestCase{
			name:       "literal percent",
			makeZettel: makeZettelWithDate,
			format:     "%%",
			output:     "%",
		},
		printfFormatterTestCase{
			name: "date then tags",
			makeZettel: func() (z *lib.Zettel) {
				z = makeZettelWithDate()
				z.Metadata.SetStringTags([]string{"some-tag"})
				return
			},
			format: "%w, %t",
			output: "2021-10-10, some-tag",
		},
		printfFormatterTestCase{
			name: "newlines",
			makeZettel: func() (z *lib.Zettel) {
				z = makeZettelWithDate()
				z.Metadata.SetStringTags([]string{"f-filename.extension"})
				return
			},
			format: "%f\n",
			output: `filename.extension
`,
		},
		printfFormatterTestCase{
			name: "body",
			makeZettel: func() (z *lib.Zettel) {
				z = makeZettelWithDate()
				z.Body = `
90210
				`
				return
			},
			format: "%b",
			output: "90210",
		},
	}
}

func TestPrintfFormat(t *testing.T) {
	for _, tc := range getPrintfTestCases(t) {
		t.Run(
			tc.name,
			func(t *testing.T) {
				formatted := FormatZettel(tc.makeZettel(), tc.format)

				if formatted != tc.output {
					t.Errorf("Formatted string was '%s', wanted '%s'", formatted, tc.output)
				}
			},
		)
	}
}
