package metadata

type ITag interface {
	Set(string) error
	Tag() string
	SearchMatchTags() TagSet
}
