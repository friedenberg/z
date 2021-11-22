package metadata

import (
	"regexp"
	"strings"
)

var (
	regexExpandTagsHyphens *regexp.Regexp
)

func init() {
	regexExpandTagsHyphens = regexp.MustCompile(`[^-]-+[^-]`)
}

type ITag interface {
	Set(string) error
	Tag() string
	SearchMatchTags() TagSet
	Match(string) bool
}

type Tag string

func (t *Tag) Set(s string) (err error) {
	*t = Tag(s)
	return
}

func (t Tag) Tag() string {
	return string(t)
}

func (t Tag) SearchMatchTags() (expanded TagSet) {
	if t == "" {
		return
	}

	hyphens := regexExpandTagsHyphens.FindAllIndex([]byte(t), -1)
	expanded = MakeTagSet()

	smt := SearchMatchTag(t)
	expanded.Add(&smt)

	if hyphens == nil {
		return
	}

	end := len(t)
	prevLocEnd := 0

	for i, loc := range hyphens {
		locStart := loc[0] + 1
		locEnd := loc[1] - 1
		t1 := SearchMatchTag(t[0:locStart])
		t2 := SearchMatchTag(t[locEnd:end])
		expanded.Add(&t1)
		expanded.Add(&t2)

		if 0 < i && i < len(hyphens) {
			t1 := SearchMatchTag(t[prevLocEnd:locStart])
			expanded.Add(&t1)
		}

		prevLocEnd = locEnd
	}

	return
}

func (t Tag) Match(q string) bool {
	return strings.Contains(string(t), q)
}
