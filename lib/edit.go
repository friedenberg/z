package lib

import (
	"os/exec"
)

func (z *Zettel) Edit(basePath string) (err error) {
	// cmd := exec.Command(editor, args...)
	// cmd.Dir = e.BasePath

	// cmd.Stdin = os.Stdin
	// cmd.Stdout = os.Stdout

	// cmd.Run()

	c := exec.Command("open", z.Path)
	c.Dir = basePath
	err = c.Run()
	return
}
