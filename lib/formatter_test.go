package lib

import (
	"testing"
)

type printFormatterTestCaseMakeZettelFunc func() *Zettel

type printfFormatterTestCase struct {
	name       string
	makeZettel printFormatterTestCaseMakeZettelFunc
	format     string
	output     string
}

func getPrintfTestCases(t *testing.T) []printfFormatterTestCase {
	k := FileStore{}
	umwelt := Umwelt{
		Kasten: Kasten{Local: &k},
	}

	makeZettelWithDate := func() (z *Zettel) {
		z = &Zettel{
			Umwelt: &umwelt,
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
			makeZettel: func() (z *Zettel) {
				z = makeZettelWithDate()
				z.Metadata.SetStringTags([]string{"some-tag"})
				return
			},
			format: "%w, %t",
			output: "2021-10-10, some-tag",
		},
		printfFormatterTestCase{
			name: "newlines",
			makeZettel: func() (z *Zettel) {
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
			makeZettel: func() (z *Zettel) {
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
				formatFunc := MakePrintfFormatFunc(tc.format)
				formatted := formatFunc(tc.makeZettel())

				if formatted != tc.output {
					t.Errorf("Formatted string was '%s', wanted '%s'", formatted, tc.output)
				}
			},
		)
	}
}
