package commands

import (
	"flag"
	"fmt"
	"os/exec"

	"github.com/friedenberg/z/lib"
)

func GetSubcommandOpen(f *flag.FlagSet) CommandRunFunc {
	var shouldPerformAction bool
	f.BoolVar(&shouldPerformAction, "action", false, "")

	return func(e Env) (err error) {
		path := f.Arg(0)

		if path == "" {
			//TODO
		}

		z := &lib.Zettel{}

		z.HydrateFromFilePath(path)

		cmd := "mvim"
		args := []string{path}

		if shouldPerformAction {
			cmd, args = GetOpenCmdAndArgs(z, path)
		}

		fmt.Println(cmd, args)
		c := exec.Command(cmd, args...)
		c.Dir = e.ZettelPath
		return c.Run()
	}
}

func GetOpenCmdAndArgs(z *lib.Zettel, path string) (cmd string, args []string) {
	switch z.Metadata.Kind {
	case "file":
		cmd = "open"
		fmt.Println(z.Metadata.File)
		args = []string{z.Metadata.File}
		//TODO check if file exists
	case "pb":
		cmd = "open"
		args = []string{z.Metadata.Url}
	case "script":
		//TODO
		cmd = path
	}

	return
}
