package filter

import (
	"strconv"
)

type TagExclusions struct {
	shouldExclude bool
}

func MakeTagExclusions() TagExclusions {
	return TagExclusions{shouldExclude: true}
}

func (t TagExclusions) String() string {
	return ""
}

func (t *TagExclusions) Set(s string) (err error) {
	v, err := strconv.ParseBool(s)

	if err != nil {
		return
	}

	t.shouldExclude = !v

	return
}

func (t TagExclusions) getNotOrFilterFromExcludedTags(excludedTags []string) (f Filter) {
	fs := make([]Filter, len(excludedTags))

	for i, f := range excludedTags {
		fs[i] = Tag(f)
	}

	f = Not(MakeOr(fs...))

	return
}

func (t TagExclusions) WithFilter(f Filter, excludedTags []string) (f1 Filter) {
	f1 = f

	if t.shouldExclude {
		f1 = MakeAnd(
			t.getNotOrFilterFromExcludedTags(excludedTags),
			f,
		)
	}

	return
}

func (t TagExclusions) IsBoolFlag() bool {
	return true
}
