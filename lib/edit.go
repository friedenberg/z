package lib

import (
	"os/exec"
)

func (z *Zettel) Edit() (err error) {
	// cmd := exec.Command(editor, args...)
	// cmd.Dir = e.BasePath

	// cmd.Stdin = os.Stdin
	// cmd.Stdout = os.Stdout

	// cmd.Run()

	c := exec.Command("open", z.Path)
	c.Dir = z.Kasten.BasePath
	err = c.Run()
	return
}
