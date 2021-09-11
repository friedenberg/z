package util

import (
	"io"
	"os"
)

type openFilesGuard struct {
	channel chan struct{}
}

var OpenFilesGuardInstance *openFilesGuard

func init() {
	OpenFilesGuardInstance = &openFilesGuard{
		channel: make(chan struct{}, 240),
	}
}

func (g *openFilesGuard) Lock() {
	g.channel <- struct{}{}
}

func (g *openFilesGuard) Unlock() {
	<-g.channel
}

func (g *openFilesGuard) Create(s string) (f *os.File, err error) {
	g.Lock()
	f, err = os.Create(s)

	if err != nil {
		g.Unlock()
	}

	return
}

func (g *openFilesGuard) OpenFile(name string, flag int, perm os.FileMode) (f *os.File, err error) {
	g.Lock()
	return os.OpenFile(name, flag, perm)
}

func (g *openFilesGuard) Open(s string) (f *os.File, err error) {
	g.Lock()
	return os.Open(s)
}

func (g *openFilesGuard) Close(f io.Closer) error {
	defer g.Unlock()
	return f.Close()
}
