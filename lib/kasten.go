package lib

import (
	"path"
	"path/filepath"
)

type Kasten struct {
	Implementation KastenImplementation
	BasePath       string
	Index          Index
}

func (e *Kasten) GetAll() (zettels []string, err error) {
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
