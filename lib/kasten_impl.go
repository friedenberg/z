package lib

type OnZettelWriteFunc func(*Zettel, error) error

type KastenImplementation interface {
	GetAll() (zs []int64, err error)
	Hydrate(readBody bool) (err error)
	ReadMetadata() (err error)
	ParseMetadata() (err error)
	ReadMetadataAndBody() (err error)
	Write(onWriteFunc OnZettelWriteFunc) (err error)
}
