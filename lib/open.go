package lib

import (
	"fmt"
	"os/exec"

	"github.com/friedenberg/z/util"
)

func (z *Zettel) Open(basePath string) (err error) {
	cmds, err := getOpenCmds(z)

	if err != nil {
		return err
	}

	for _, c := range cmds {
		c.Dir = basePath
		err = c.Run()

		if err != nil {
			return
		}
	}

	return
}

func getOpenCmds(z *Zettel) (c []*exec.Cmd, err error) {
	if z.HasFile() {
		if !util.FileExists(z.IndexData.File) {
			err = fmt.Errorf("%s: file does not exist", z.IndexData.File)
			return
		}

		c = append(c, exec.Command("open", z.IndexData.File))
	}

	if z.HasUrl() {
		c = append(c, exec.Command("open", z.IndexData.Url))
	}

	return
}
