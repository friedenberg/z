package lib

import (
	"testing"
	"time"
)

type printFormatterTestCaseMakeZettelFunc func() *Zettel

type printfFormatterTestCase struct {
	name       string
	makeZettel printFormatterTestCaseMakeZettelFunc
	format     string
	output     string
}

func getPrintfTestCases(t *testing.T) []printfFormatterTestCase {
	kasten := FilesAndGit{}

	makeZettelWithDate := func() (z *Zettel) {
		time, err := time.Parse("2006-01-02", "2021-07-26")

		if err != nil {
			t.Fatal(err)
		}

		z = &Zettel{
			FilesAndGit: &kasten,
		}

		z.InitFromTime(time)

		return
	}

	return []printfFormatterTestCase{
		printfFormatterTestCase{
			name:       "just date",
			makeZettel: makeZettelWithDate,
			format:     "%w",
			output:     "2021-07-26",
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
				z.Metadata.Tags = []string{"some-tag"}
				return
			},
			format: "%w, %t",
			output: "2021-07-26, some-tag",
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
