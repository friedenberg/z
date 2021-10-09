package metadata

import (
	"regexp"
)

var (
	regexExpandTagsHyphens *regexp.Regexp
)

func init() {
	regexExpandTagsHyphens = regexp.MustCompile(`[^-]-+[^-]`)
}

func ExpandTagsAsTags(t string) (expanded []ITag) {
	if t == "" {
		return
	}

	hyphens := regexExpandTagsHyphens.FindAllIndex([]byte(t), -1)
	it := Tag(t)
	expanded = []ITag{&it}

	if hyphens == nil {
		return
	}

	end := len(t)
	prevLocEnd := 0

	for i, loc := range hyphens {
		locStart := loc[0] + 1
		locEnd := loc[1] - 1
		t1 := Tag(t[0:locStart])
		t2 := Tag(t[locEnd:end])
		expanded = append(expanded, &t1)
		expanded = append(expanded, &t2)

		if 0 < i && i < len(hyphens) {
			t1 := Tag(t[prevLocEnd:locStart])
			expanded = append(expanded, &t1)
		}

		prevLocEnd = locEnd
	}

	return
}

func ExpandTags(t string) (expanded []string) {
	e := ExpandTagsAsTags(t)
	expanded = make([]string, len(e))

	for i, t := range e {
		expanded[i] = string(t.Tag())
	}

	return
}
