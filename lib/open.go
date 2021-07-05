package lib

import (
	"fmt"
	"os/exec"

	"github.com/friedenberg/z/util"
)

func (z *Zettel) Open(basePath string) (err error) {
	c, err := getOpenCmd(z)

	if err != nil {
		return err
	}

	c.Dir = basePath
	err = c.Run()
	return
}

func getOpenCmd(z *Zettel) (c *exec.Cmd, err error) {
	switch z.IndexData.Kind {
	case "file":
		if !util.FileExists(z.IndexData.File) {
			err = fmt.Errorf("%s: file does not exist", z.IndexData.File)
			return
		}

		c = exec.Command("open", z.IndexData.File)

	case "pb":
		c = exec.Command("open", z.IndexData.Url)
	}

	return
}
