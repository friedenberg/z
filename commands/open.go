package commands

import (
	"flag"
	"fmt"
	"os/exec"

	"github.com/friedenberg/z/lib"
	"github.com/friedenberg/z/util"
)

func GetSubcommandOpen(f *flag.FlagSet) CommandRunFunc {
	var shouldEdit bool
	var shouldOpen bool

	f.BoolVar(&shouldEdit, "edit", true, "")
	f.BoolVar(&shouldOpen, "action", false, "")

	return func(e Env) (err error) {
		processor := MakeProcessor(
			e,
			f.Args(),
			&NullPutter{Channel: make(PutterChannel)},
		)

		processor.actioner = func(i int, z *lib.Zettel) (err error) {
			var c *exec.Cmd

			if shouldEdit {
				c := exec.Command("open", z.Path)
				c.Dir = e.ZettelPath
				err = c.Run()

				if err != nil {
					return
				}
			}

			if shouldOpen {
				c, err = getOpenCmd(z)

				if err != nil {
					return err
				}

				c.Dir = e.ZettelPath
				err = c.Run()
			}

			return
		}

		err = processor.Run()

		return
	}
}

func getOpenCmd(z *lib.Zettel) (c *exec.Cmd, err error) {
	switch z.Metadata.Kind {
	case "file":
		if !util.FileExists(z.Metadata.File) {
			err = fmt.Errorf("%s: file does not exist", z.Metadata.File)
			return
		}

		c = exec.Command("open", z.Metadata.File)

	case "pb":
		c = exec.Command("open", z.Metadata.Url)
	}

	return
}
