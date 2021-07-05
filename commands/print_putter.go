package commands

import (
	"fmt"

	"github.com/friedenberg/z/lib"
)

type StdoutPutter struct {
	channel          PutterChannel
	shouldPrintComma bool
}

func MakePutter() (putter Putter) {
	putter = &StdoutPutter{channel: make(PutterChannel)}
	return
}

func (p *StdoutPutter) GetChannel() PutterChannel {
	return p.channel
}

func (p *StdoutPutter) Print() {
	fmt.Print(`{"items":[`)
	defer fmt.Print(`]}`)

	for z := range p.channel {
		p.PrintZettel(z)
	}
}

func (p *StdoutPutter) PrintZettel(z *lib.Zettel) {
	if p.shouldPrintComma {
		fmt.Print(",")
	}

	fmt.Print(z.AlfredData.Json)

	p.shouldPrintComma = true
}
