package lib

import "os/exec"

func (z *Zettel) Edit(basePath string) (err error) {
	c := exec.Command("open", z.Path)
	c.Dir = basePath
	err = c.Run()
	return
}
