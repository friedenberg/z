package stdprinter

import (
	"fmt"
	"os"
	"sync"
)

type printerLine struct {
	file    *os.File
	line    string
	isDebug bool
}

var (
	printerChannel   chan printerLine
	waitGroup        *sync.WaitGroup
	shouldPrintDebug bool
)

func init() {
	printerChannel = make(chan printerLine)
	waitGroup = &sync.WaitGroup{}
	waitGroup.Add(1)

	shouldPrintDebug = false

	go func() {
		defer waitGroup.Done()

		for printerLine := range printerChannel {
			if printerLine.isDebug && !shouldPrintDebug {
				continue
			}

			fmt.Fprint(printerLine.file, printerLine.line)
		}
	}()
}

func WaitForPrinter() {
	close(printerChannel)
	waitGroup.Wait()
}

func SetDebug(d bool) {
	shouldPrintDebug = d
}
