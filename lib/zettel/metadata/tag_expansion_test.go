package metadata

import (
	"testing"
)

type ExpandTagsTestCase struct {
	description  string
	tag          string
	expandedTags []string
}

func GetTestExpandTagsTestCases(t *testing.T) []ExpandTagsTestCase {
	return []ExpandTagsTestCase{
		ExpandTagsTestCase{
			"empty",
			"",
			[]string{},
		},
		ExpandTagsTestCase{
			"no hyphens",
			"tag",
			[]string{"tag"},
		},
		ExpandTagsTestCase{
			"one hyphen",
			"tag-part",
			[]string{
				"part",
				"tag",
				"tag-part",
			},
		},
		ExpandTagsTestCase{
			"two hyphens",
			"p-2021-zettel",
			[]string{
				"2021",
				"2021-zettel",
				"p",
				"p-2021",
				"p-2021-zettel",
				"zettel",
			},
		},
	}
}

func TestExpandTags(t *testing.T) {
	for _, testCase := range GetTestExpandTagsTestCases(t) {
		t.Run(
			testCase.description,
			func(t *testing.T) {
				tag := Tag(testCase.tag)
				actual := tag.SearchMatchTags().Strings()

				if len(testCase.expandedTags) != len(actual) {
					t.Errorf(
						"Expanded tags was '%q', wanted '%q'",
						actual,
						testCase.expandedTags,
					)

					return
				}

				for i, a := range actual {
					e := testCase.expandedTags[i]
					if a != e {
						t.Errorf(
							"Expanded tags was '%q' at %d, wanted '%q'",
							a,
							i,
							e,
						)
					}
				}
			},
		)
	}
}
