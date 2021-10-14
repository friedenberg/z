package lib

import (
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/friedenberg/z/lib/zettel"
	"github.com/friedenberg/z/util"
)

//TODO-P2 move to lib/kasten
type FileStore struct {
	basePath string
}

func (s FileStore) BasePath() string {
	return s.basePath
}

func (e FileStore) GetAll() (zettels []string, err error) {
	glob := filepath.Join(e.BasePath(), "*.md")
	zettels, err = filepath.Glob(glob)
	return
}

func (e FileStore) GetNormalizedPath(a string) (b string, err error) {
	if filepath.IsAbs(a) {
		b = a
	} else {
		b, err = filepath.Abs(path.Join(e.BasePath(), a))
	}

	return
}

func (k *FileStore) NewId() (id zettel.Id, err error) {
	currentTime := time.Now()

	for {
		p := MakePathFromTime(k.BasePath(), currentTime)

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

func (k FileStore) CommitTransaction(u Umwelt) (err error) {
	readAndWrite := func(z *Zettel) (err error) {
		err = z.Hydrate(true)

		if os.IsNotExist(err) {
			err = nil
		} else if err != nil {
			return
		}

		err = z.Write(nil)

		if err != nil {
			return
		}

		return
	}

	for _, z := range u.Transaction.Added() {
		err = readAndWrite(z)

		if err != nil {
			return
		}
	}

	for _, z := range u.Transaction.Modified() {
		err = readAndWrite(z)

		if err != nil {
			return
		}
	}

	return
}