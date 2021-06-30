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
		processor := MakeProcessor(
			e,
			f.Args(),
			&NullPutter{Channel: make(PutterChannel)},
		)

		processor.parallelAction = func(i int, z *lib.Zettel) (err error) {
			cmd := "mvim"
			args := []string{z.Path}

			if shouldPerformAction {
				cmd, args = GetOpenCmdAndArgs(z, z.Path)
			}

			c := exec.Command(cmd, args...)
			c.Dir = e.ZettelPath
			return c.Run()
		}

		err = processor.Run()

		return
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
