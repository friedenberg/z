package lib

type OnZettelWriteFunc func(*KastenZettel, error) error

type KastenImplementation interface {
	GetAll() (zs []int64, err error)
	NewId() (int64, error)
	Hydrate(readBody bool) (err error)
	ReadMetadata() (err error)
	ParseMetadata() (err error)
	ReadMetadataAndBody() (err error)
	Write(onWriteFunc OnZettelWriteFunc) (err error)
}
