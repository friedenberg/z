package kasten

type RemoteImplementation interface {
	InitFromOptions(map[string]interface{}) (err error)
	//AddFile(p string) error
}

type LocalImplementation interface {
	RemoteImplementation
}
