package commands

import (
	"os/user"
	"path"
	"path/filepath"
)

type Env struct {
	ZettelPath string
}

func GetDefaultEnv() (e Env, err error) {
	usr, err := user.Current()

	if err != nil {
		return
	}

	home := usr.HomeDir
	e = Env{
		ZettelPath: path.Join(home, "Zettelkasten"),
	}

	return
}

func (e Env) GetAllZettels() (zettels []string, err error) {
	glob := filepath.Join(e.ZettelPath, "*.md")
	zettels, err = filepath.Glob(glob)
	return
}

func (e Env) GetNormalizedPath(a string) (b string, err error) {
	if filepath.IsAbs(a) {
		b = a
	} else {
		b, err = filepath.Abs(path.Join(e.ZettelPath, a))
	}

	return
}
