package kasten

import "github.com/friedenberg/z/lib/zettel"

type RemoteImplementation interface {
	InitFromOptions(map[string]interface{}) (err error)
	CopyFileTo(localPath string, fd zettel.FileDescriptor) (err error)
	CopyFileFrom(localPath string, fd zettel.FileDescriptor) (err error)
}

type LocalImplementation interface {
	RemoteImplementation
}
