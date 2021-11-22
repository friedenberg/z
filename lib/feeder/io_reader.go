package feeder

import (
	"bufio"
	"io"
)

type IoReader struct {
	stream  io.Reader
	argChan chan string
}

func MakeIoReader(stream io.Reader) IoReader {
	return IoReader{
		stream:  stream,
		argChan: make(chan string, 0),
	}
}

func (f IoReader) GetChan() <-chan string {
	return f.argChan
}

func (f IoReader) Run() {
	defer close(f.argChan)

	scanner := bufio.NewScanner(f.stream)

	for scanner.Scan() {
		f.argChan <- scanner.Text()
	}
}
