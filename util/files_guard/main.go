package files_guard

import (
	"io"
	"os"
)

type openFilesGuard struct {
	channel chan struct{}
}

var openFilesGuardInstance *openFilesGuard

func init() {
	openFilesGuardInstance = &openFilesGuard{
		//TODO read from OS and determine dynamically
		channel: make(chan struct{}, 240),
	}
}

func (g *openFilesGuard) Lock() {
	g.channel <- struct{}{}
}

func (g *openFilesGuard) Unlock() {
	<-g.channel
}

func Create(s string) (f *os.File, err error) {
	openFilesGuardInstance.Lock()
	f, err = os.Create(s)

	if err != nil {
		openFilesGuardInstance.Unlock()
	}

	return
}

func OpenFile(name string, flag int, perm os.FileMode) (f *os.File, err error) {
	openFilesGuardInstance.Lock()
	return os.OpenFile(name, flag, perm)
}

func Open(s string) (f *os.File, err error) {
	openFilesGuardInstance.Lock()
	return os.Open(s)
}

func Close(f io.Closer) error {
	defer openFilesGuardInstance.Unlock()
	return f.Close()
}
