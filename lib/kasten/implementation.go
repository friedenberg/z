package kasten

type RemoteImplementation interface {
	InitFromOptions(map[string]interface{}) (err error)
	CopyFileTo(localPath string, h RemoteFileHandle) (err error)
	CopyFileFrom(localPath string, h RemoteFileHandle) (err error)
}

type LocalImplementation interface {
	RemoteImplementation
}
