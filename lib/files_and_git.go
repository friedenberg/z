package lib

import (
	"path"
	"path/filepath"
	"time"

	"github.com/friedenberg/z/lib/zettel"
	"github.com/friedenberg/z/util"
)

//TODO-P2 move to lib/kasten
type FilesAndGit struct {
	BasePath        string
	GitEnabled      bool
	GitAnnexEnabled bool
}

func (e *FilesAndGit) GetAll() (zettels []string, err error) {
	glob := filepath.Join(e.BasePath, "*.md")
	zettels, err = filepath.Glob(glob)
	return
}

func (e *FilesAndGit) GetNormalizedPath(a string) (b string, err error) {
	if filepath.IsAbs(a) {
		b = a
	} else {
		b, err = filepath.Abs(path.Join(e.BasePath, a))
	}

	return
}

func (k *FilesAndGit) NewId() (id zettel.Id, err error) {
	currentTime := time.Now()

	for {
		p := MakePathFromTime(k.BasePath, currentTime)

		if util.FileExists(p) {
			d, err := time.ParseDuration("1s")

			if err != nil {
				panic(err)
			}

			currentTime = currentTime.Add(d)
		} else {
			id = zettel.Id(currentTime.Unix())
			return
		}
	}

	return
}
