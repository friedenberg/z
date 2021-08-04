package lib

import (
	"fmt"
	"os/exec"

	"github.com/friedenberg/z/util"
)

func (z *Zettel) Open() (err error) {
	cmds, err := getOpenCmds(z)

	if err != nil {
		return err
	}

	for _, c := range cmds {
		c.Dir = z.Kasten.BasePath
		err = c.Run()

		if err != nil {
			err = fmt.Errorf("opening attachment %s: %w", c, err)
			return
		}
	}

	return
}

func getOpenCmds(z *Zettel) (c []*exec.Cmd, err error) {
	if z.HasFile() {
		if !util.FileExists(z.FilePath()) {
			err = fmt.Errorf("%s: file does not exist", z.FilePath())
			return
		}

		c = append(c, exec.Command("open", z.FilePath()))
	}

	if z.HasUrl() {
		c = append(c, exec.Command("open", z.IndexData.Url))
	}

	return
}
