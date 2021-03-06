package metadata

import (
	"testing"
)

func TestMetadataToJSON(t *testing.T) {
	assertMetadataSetStringTags(
		t,
		[]string{"t-test-tag"},
		`["t-test-tag"]`,
	)
}

func TestMetadataToJSON2(t *testing.T) {
	assertMetadataSetStringTags(
		t,
		[]string{},
		`[]`,
	)
}

func TestMetadataToJSON4(t *testing.T) {
	assertMetadataSetStringTags(
		t,
		[]string{"ok", "wow", "u-google.com/thispath"},
		`["ok","u-https://google.com/thispath","wow"]`,
	)
}

func TestMetadataToJSON5(t *testing.T) {
	assertMetadataUnmarshalJSON(
		t,
		`["ok","wow","u-https://google.com/thispath"]`,
		`["ok","u-https://google.com/thispath","wow"]`,
	)
}

func TestMetadataToJSON6(t *testing.T) {
	assertMetadataUnmarshalJSON(
		t,
		`["ok","u-https://google.com/thispath","wow"]`,
		`["ok","u-https://google.com/thispath","wow"]`,
	)
}

func TestMetadataToJSON7(t *testing.T) {
	assertMetadataUnmarshalJSON(
		t,
		`["ok","u-https://google.com/thispath","wow"]`,
		`["ok","u-https://google.com/thispath","wow"]`,
	)
}

func TestMetadataToJSON8(t *testing.T) {
	assertMetadataSetStringTags(
		t,
		[]string{""},
		`[]`,
	)
}

func TestMetadataToJSONNewFile(t *testing.T) {
	assertMetadataSetStringTags(
		t,
		[]string{"nf-the_path"},
		`["nf-the_path"]`,
	)

	assertMetadataSetStringTagsNewFileTag(
		t,
		[]string{"nf-the_path"},
		"the_path",
	)
}

func assertMetadataUnmarshalJSON(t *testing.T, input, expected string) {
	t.Helper()

	var m Metadata
	err := m.UnmarshalJSON([]byte(input))

	if err != nil {
		t.Fatalf("failed to set from tag: %s", err)
	}

	b, err := m.MarshalJSON()

	if err != nil {
		t.Fatalf("failed to marshal json: %s", err)
	}

	json := string(b)

	if expected != json {
		t.Errorf("Actual json was '%s', wanted '%s'", json, expected)
	}
}

func assertMetadataSetStringTags(t *testing.T, in []string, expected string) {
	t.Helper()
	var m Metadata
	err := m.SetStringTags(in)

	if err != nil {
		t.Fatalf("failed to set string tags: %s", err)
	}

	b, err := m.MarshalJSON()

	if err != nil {
		t.Fatalf("failed to marshal json: %s", err)
	}

	json := string(b)

	if expected != json {
		t.Errorf("\n  actual: '%s'\nexpected: '%s'", json, expected)
	}
}

func assertMetadataSetStringTagsNewFileTag(t *testing.T, in []string, expected string) {
	t.Helper()
	var m Metadata
	err := m.SetStringTags(in)

	if err != nil {
		t.Fatalf("failed to set string tags: %s", err)
	}

	tag, ok := m.NewFile()

	if !ok {
		t.Fatalf("metadata missing new file in metadata:\n%#v", m)
	}

	actual := tag.Path

	if expected != actual {
		t.Errorf("\n  actual: '%s'\nexpected: '%s'", actual, expected)
	}
}

func TestMetadataToJSON9(t *testing.T) {
	assertYamlMatches(t,
		`---
- ok
- u-https://www.google.com/thispath
- t-wow
...
`,
		`---
- ok
- u-https://www.google.com/thispath
- t-wow
...
`,
	)
}

func assertYamlMatches(t *testing.T, i string, e string) {
	t.Helper()

	var m Metadata
	err := m.Set(i)

	if err != nil {
		t.Fatalf("failed to set from tag: %s", err)
	}

	b, err := m.ToYAMLWithBoundary()

	if err != nil {
		t.Errorf("failed to marshal json: %s", err)
	}

	if b != e {
		t.Errorf("Actual json was '%s', wanted '%s'", b, e)
	}
}
