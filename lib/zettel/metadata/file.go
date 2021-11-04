package metadata

type File interface {
	ITag
	Extension() string
	FilePath(string) string
}
