package commands

import (
	"github.com/friedenberg/z/lib"
)

type PutterChannel chan *lib.Zettel

type Putter interface {
	GetChannel() PutterChannel
	Print()
}

func MakeNullPutter() Putter {
	return &NullPutter{
		Channel: make(PutterChannel),
	}
}

type NullPutter struct {
	Channel PutterChannel
}

func (p *NullPutter) GetChannel() PutterChannel {
	return p.Channel
}

func (p *NullPutter) Print() {
	for _ = range p.Channel {
	}
}
