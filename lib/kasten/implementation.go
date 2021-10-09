package kasten

import "github.com/friedenberg/z/lib/zettel/metadata"

type RemoteImplementation interface {
	InitFromOptions(map[string]interface{}) (err error)
	CopyFileTo(localPath string, fd metadata.FileDescriptor) (err error)
	CopyFileFrom(localPath string, fd metadata.FileDescriptor) (err error)
}

type LocalImplementation interface {
	RemoteImplementation
}
