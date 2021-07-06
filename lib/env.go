package lib

import (
	"os/user"
	"path"
	"path/filepath"
	"sync"
)

type Env struct {
	BasePath string
	ZettelPool
}

func GetDefaultEnv() (e *Env, err error) {
	usr, err := user.Current()

	if err != nil {
		return
	}

	e = &Env{
		BasePath: path.Join(usr.HomeDir, "Zettelkasten"),
	}

	e.ZettelPool = &zettelPool{
		Pool: sync.Pool{
			New: func() interface{} {
				z := new(Zettel)
				z.Env = e
				return z
			},
		},
		env: e,
	}

	return
}

func (e *Env) GetAllZettels() (zettels []string, err error) {
	glob := filepath.Join(e.BasePath, "*.md")
	zettels, err = filepath.Glob(glob)
	return
}

func (e *Env) GetNormalizedPath(a string) (b string, err error) {
	if filepath.IsAbs(a) {
		b = a
	} else {
		b, err = filepath.Abs(path.Join(e.BasePath, a))
	}

	return
}
