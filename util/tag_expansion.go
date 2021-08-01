package util

import (
	"regexp"
)

var (
	regexExpandTagsHyphens *regexp.Regexp
)

func init() {
	regexExpandTagsHyphens = regexp.MustCompile(`[^-]-+[^-]`)
}

func ExpandTags(t string) (expanded []string) {
	hyphens := regexExpandTagsHyphens.FindAllIndex([]byte(t), -1)

	if hyphens == nil {
		return
	}

	end := len(t)
	prevLocEnd := 0

	for i, loc := range hyphens {
		locStart := loc[0] + 1
		locEnd := loc[1] - 1
		expanded = append(expanded, t[0:locStart])
		expanded = append(expanded, t[locEnd:end])

		if 0 < i && i < len(hyphens) {
			expanded = append(expanded, t[prevLocEnd:locStart])
		}

		prevLocEnd = locEnd
	}

	return
}
