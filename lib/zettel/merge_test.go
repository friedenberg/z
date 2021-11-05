package zettel

import (
	"bytes"
	"strings"
	"testing"

	"github.com/friedenberg/z/lib/zettel/metadata"
)

func TestZettelMerge1(t *testing.T) {
	assertZettelMerge(
		t,
		`---
- title
- tag-1
...

body content
`,
		`---
- secondtitle
- tag-2
...

second body content
`,
		`---
- title secondtitle
- tag-1
- tag-2
...

body content

second body content
`,
	)
}

func assertZettelMerge(t *testing.T, z1Yaml, z2Yaml, expected string) {
	t.Helper()
	var err error

	z1 := &Zettel{}
	z2 := &Zettel{}

	read := func(z *Zettel, y string) {
		z.Note = Note{
			Metadata: metadata.MakeMetadata(),
		}

		r := strings.NewReader(y)
		err = z.ReadFrom(r, true)

		if err != nil {
			t.Fatalf("failed to read yaml from input: %s", err)
		}
	}

	read(z1, z1Yaml)
	read(z2, z2Yaml)

	err = z1.Merge(z2)

	if err != nil {
		t.Fatalf("failed to merge: %s", err)
	}

	sb := new(bytes.Buffer)
	err = z1.WriteTo(sb)

	if err != nil {
		t.Fatalf("failed to read merged yaml: %s", err)
	}

	actual := sb.String()

	if actual != expected {
		t.Errorf("\nexpected: '%q'\nactual:   '%q'", expected, actual)
	}
}
