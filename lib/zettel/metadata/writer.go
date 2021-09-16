package metadata

import (
	"strings"

	"gopkg.in/yaml.v2"
)

func (m Metadata) ToSortedTags() (t []string) {
	t = make([]string, 0, m.stringTags.Len())

	if m.description != "" {
		t = append(t, m.description)
	}

	if u, ok := m.Url(); ok {
		t = append(t, u.Tag())
	}

	//TODO-P4 move to cleanup
	for _, f := range m.stringTags.Tags() {
		if _, ok := f.(*SearchMatchTag); ok {
			continue
		}

		t = append(t, f.Tag())
	}

	for _, f := range m.Files() {
		t = append(t, f.Tag())
	}

	return
}

func (m Metadata) ToYAML() (h string, err error) {
	t := m.ToSortedTags()
	b, err := yaml.Marshal(t)

	if err != nil {
		return
	}

	h = string(b)
	return
}

func (m Metadata) ToYAMLWithBoundary() (h string, err error) {
	w := &strings.Builder{}

	w.WriteString(MetadataStartSequence)

	y, err := m.ToYAML()

	if err != nil {
		return
	}

	if y != "[]\n" {
		w.WriteString(y)
	}

	if err != nil {
		return
	}

	w.WriteString(MetadataEndSequence)

	h = w.String()

	return
}
