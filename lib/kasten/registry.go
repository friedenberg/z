package kasten

import "golang.org/x/xerrors"

type factory func() RemoteImplementation

var (
	registryInstance map[string]factory
)

func init() {
	registryInstance = make(map[string]factory)
	Register("files", func() RemoteImplementation { return &Files{} })
}

func Register(n string, f factory) (err error) {
	if _, ok := registryInstance[n]; ok {
		err = xerrors.Errorf("Multiple implementations with name: '%s'", n)
		return
	}

	registryInstance[n] = f

	return
}

func GetLocal(n string) (i LocalImplementation, ok bool) {
	ir, ok := GetRemote(n)
	i, ok = ir.(LocalImplementation)

	return
}

func GetRemote(n string) (i RemoteImplementation, ok bool) {
	//TODO normalize n
	f, ok := registryInstance[n]
	i = f()
	return
}
