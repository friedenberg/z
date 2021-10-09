package metadata

type ITag interface {
	Set(string) error
	Tag() string
}

type Tag string

func (t *Tag) Set(s string) (err error) {
	*t = Tag(s)
	return
}

func (t Tag) Tag() string {
	return string(t)
}
