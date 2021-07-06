package util

import "sync"

var OpenFilesGuardInstance sync.Locker

func init() {
	OpenFilesGuardInstance = &openFilesGuard{
		channel: make(chan struct{}, 240),
	}
}

type openFilesGuard struct {
	channel chan struct{}
}

func (g *openFilesGuard) Lock() {
	g.channel <- struct{}{}
}

func (g *openFilesGuard) Unlock() {
	<-g.channel
}
