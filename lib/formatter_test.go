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
	env := Env{}

	makeZettelWithDate := func() (z *Zettel) {
		time, err := time.Parse("2006-01-02", "2021-07-26")

		if err != nil {
			t.Fatal(err)
		}

		z = &Zettel{
			Env: &env,
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
			name:       "literal percent",
			makeZettel: makeZettelWithDate,
			format:     "%w,a:%a,p:%r,t:%t",
			output:     "2021-07-26,a:,p:,t:",
		},
		printfFormatterTestCase{
			name: "literal percent",
			makeZettel: func() (z *Zettel) {
				z = makeZettelWithDate()
				z.IndexData.Areas = []string{"some-area"}
				z.IndexData.Projects = []string{"some-project"}
				z.IndexData.Tags = []string{"some-tag"}
				return
			},
			format: "%w,a:%a,p:%r,t:%t",
			output: "2021-07-26,a:some-area,p:some-project,t:some-tag",
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
