package stdprinter

import (
	"fmt"
	"os"
)

func Error(err error) {
	printerChannel <- printerLine{
		file: os.Stderr,
		line: fmt.Sprintf("%+v", err),
	}
}

func Errf(f string, a ...interface{}) {
	printerChannel <- printerLine{
		file: os.Stderr,
		line: fmt.Sprintf(f, a...),
	}
}

func Err(a ...interface{}) {
	printerChannel <- printerLine{
		file: os.Stderr,
		line: fmt.Sprintln(a...),
	}
}
