package commands

import (
	"os/user"
	"path"
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
