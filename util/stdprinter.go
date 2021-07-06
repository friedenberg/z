package util

import (
	"fmt"
	"os"
	"sync"
)

type printerLine struct {
	file *os.File
	line string
}

var (
	printerChannel      chan printerLine
	stdPrinterWaitGroup *sync.WaitGroup
)

func init() {
	printerChannel = make(chan printerLine)
	stdPrinterWaitGroup = &sync.WaitGroup{}
	stdPrinterWaitGroup.Add(1)

	go func() {
		defer stdPrinterWaitGroup.Done()

		for printerLine := range printerChannel {
			fmt.Fprint(printerLine.file, printerLine.line)
		}
	}()
}

func StdPrinterErrf(f string, a ...interface{}) {
	printerChannel <- printerLine{
		file: os.Stderr,
		line: fmt.Sprintf(f, a...),
	}
}

func StdPrinterErr(a ...interface{}) {
	printerChannel <- printerLine{
		file: os.Stderr,
		line: fmt.Sprintln(a...),
	}
}

func StdPrinterOutf(f string, a ...interface{}) {
	printerChannel <- printerLine{
		file: os.Stdout,
		line: fmt.Sprintf(f, a...),
	}
}

func StdPrinterOut(a ...interface{}) {
	printerChannel <- printerLine{
		file: os.Stdout,
		line: fmt.Sprintln(a...),
	}
}

func WaitForPrinter() {
	close(printerChannel)
	stdPrinterWaitGroup.Wait()
}
