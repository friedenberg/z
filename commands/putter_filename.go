package commands

import "fmt"

type FilenamePutter struct {
	Channel PutterChannel
}

func MakeFilenamePutter() Putter {
	return &FilenamePutter{
		Channel: make(PutterChannel),
	}
}

func (p *FilenamePutter) GetChannel() PutterChannel {
	return p.Channel
}

func (p *FilenamePutter) Print() {
	for z := range p.Channel {
		fmt.Println(z.Path)
	}
}
