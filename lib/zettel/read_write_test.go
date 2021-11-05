package zettel

import (
	"bytes"
	"strings"
	"testing"

	"github.com/friedenberg/z/lib/zettel/metadata"
)

func TestZettelRead(t *testing.T) {
	assertZettelRead(
		t,
		`---
- description
- tag-1
...

this is the body
`,
	)
}

func assertZettelRead(t *testing.T, expected string) {
	t.Helper()
	var err error

	z := &Zettel{
		Note: Note{
			Metadata: metadata.MakeMetadata(),
		},
	}

	err = z.ReadFrom(strings.NewReader(expected), true)

	if err != nil {
		t.Errorf("failed to read yaml from input: %s", err)
	}

	w := new(bytes.Buffer)
	err = z.WriteTo(w)

	if err != nil {
		t.Errorf("failed to read yaml from input: %s", err)
	}

	actual1 := w.String()

	if actual1 != expected {
		t.Errorf("\nexpected: '%q'\nactual:   '%q'", expected, actual1)
	}
}
