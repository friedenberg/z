package lib

import (
	"os/user"
	"path"
	"path/filepath"
)

type Kasten struct {
	BasePath string
	Index    Index
}

func GetDefaultKasten() (e *Kasten, err error) {
	usr, err := user.Current()

	if err != nil {
		return
	}

	e = &Kasten{
		BasePath: path.Join(usr.HomeDir, "Zettelkasten"),
		Index:    MakeIndex(),
	}

	return
}

func (e *Kasten) GetAllZettels() (zettels []string, err error) {
	glob := filepath.Join(e.BasePath, "*.md")
	zettels, err = filepath.Glob(glob)
	return
}

func (e *Kasten) GetNormalizedPath(a string) (b string, err error) {
	if filepath.IsAbs(a) {
		b = a
	} else {
		b, err = filepath.Abs(path.Join(e.BasePath, a))
	}

	return
}
