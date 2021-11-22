package metadata

type SearchMatchTag Tag

func (t SearchMatchTag) Tag() string {
	t1 := Tag(t)
	return t1.Tag()
}

func (t *SearchMatchTag) Set(s string) (err error) {
	t1 := Tag(*t)
	err = t1.Set(s)
	return
}

func (t SearchMatchTag) SearchMatchTags() (a TagSet) {
	return
}
