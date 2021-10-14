package stdprinter

import (
	"fmt"
	"os"
)

func Debug(err error) {
	printerChannel <- printerLine{
		file:    os.Stderr,
		line:    fmt.Sprintf("%+v", err),
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
