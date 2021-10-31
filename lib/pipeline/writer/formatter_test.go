package writer

import (
	"testing"

	"github.com/friedenberg/z/lib"
)

func TestFormatZettelJustDate(t *testing.T) {
	assertFormatZettel(
		t,
		nil,
		"%w",
		"2021-10-10",
	)
}

func TestFormatZettelJustConstantString(t *testing.T) {
	assertFormatZettel(
		t,
		nil,
		"wow",
		"wow",
	)
}

func TestFormatZettelJustConstantOneCharString(t *testing.T) {
	assertFormatZettel(
		t,
		nil,
		"w",
		"w",
	)
}

func TestFormatZettelLiteralPercent(t *testing.T) {
	assertFormatZettel(
		t,
		nil,
		"%%",
		"%",
	)
}

func TestFormatZettelDateThenTags(t *testing.T) {
	z := makeZettel()
	err := z.Metadata.SetStringTags([]string{"", "some-tag"})

	if err != nil {
		t.Errorf("failed to set string tags: %w", err)
	}

	assertFormatZettel(
		t,
		z,
		"%w, %t",
		"2021-10-10, some-tag",
	)
}

func TestFormatZettelNewlines(t *testing.T) {
	z := makeZettel()
	err := z.Metadata.SetStringTags([]string{"", "f-filename.extension"})

	if err != nil {
		t.Errorf("failed to set string tags: %w", err)
	}

	assertFormatZettel(
		t,
		z,
		"%f\n",
		`filenam.extension
`,
	)
}

func TestFormatZettelBody(t *testing.T) {
	z := makeZettel()
	z.Body = `
90210
				`

	assertFormatZettel(
		t,
		z,
		"%b",
		"90210",
	)
}

func makeZettel() (z *lib.Zettel) {
	k := lib.FileStore{}

	umwelt := lib.Umwelt{
		Kasten: &k,
	}

	z = &lib.Zettel{
		Umwelt: umwelt,
		Path:   "1633902356.md",
	}

	return
}

func assertFormatZettel(t *testing.T, z *lib.Zettel, format string, expected string) {
	t.Helper()

	if z == nil {
		z = makeZettel()
	}

	formatted := FormatZettel(z, format)

	if formatted != expected {
		t.Errorf("\n  actual: '%s'\nexpected: '%s'", formatted, expected)
	}
}
