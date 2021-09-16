package stdprinter

import (
	"fmt"
	"os"
)

func Debug(a ...interface{}) {
	printerChannel <- printerLine{
		file:    os.Stderr,
		line:    fmt.Sprintln(a...),
		isDebug: true,
	}
}

func Debugf(f string, a ...interface{}) {
	printerChannel <- printerLine{
		file:    os.Stderr,
		line:    fmt.Sprintf(f, a...),
		isDebug: true,
	}
}
