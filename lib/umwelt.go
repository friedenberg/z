package lib

import (
	"os"
	"path"

	"github.com/friedenberg/z/lib/feeder"
	"github.com/friedenberg/z/util/files_guard"
	"github.com/friedenberg/z/util/stdprinter"
)

type Umwelt struct {
	Kasten                 Kasten
	Index                  *Index
	BasePath               string
	Config                 Config
	TagsForNewZettels      []string
	TagsForExcludedZettels []string
	*Transaction
}

func MakeUmwelt(c Config, wd string) (k *Umwelt, err error) {
	k = &Umwelt{}
	k.Config = c

	k.BasePath = wd
	k.Index = MakeIndex()

	//find all caches
	//remove any older-version caches
	//try to load from correct version cache
	//or exit if it doesn't exist
	err = k.LoadIndexFromCache()

	if err != nil && !os.IsNotExist(err) {
		return
	}

	k.Transaction = MakeTransaction()

	return
}

func (e Umwelt) GetIndexPath() string {
	return path.Join(e.BasePath, ".zettel-cache")
}

func (u Umwelt) GetAll() feeder.Feeder {
	return feeder.MakeFeeder(
		func(c chan<- string) {
			for id, _ := range u.Index.Zettels {
				c <- id.String()
			}
		},
	)
}

func (u Umwelt) LoadIndexFromCache() (err error) {
	f, err := files_guard.Open(u.GetIndexPath())
	defer files_guard.Close(f)

	if err != nil && os.IsNotExist(err) {
		err = nil
		return
	}

	if err != nil {
		return
	}

	err = u.Index.Read(f)

	if err != nil {
		return
	}

	return
}

//TODO-P2 add lock like git: git add: exit status 128: fatal: Unable to create '/Users/sasha/Zettelkasten/.git/index.lock': File exists.
func (e Umwelt) CacheIndex() (err error) {
	f, err := files_guard.Create((e.GetIndexPath()))
	defer files_guard.Close(f)

	if err != nil {
		return
	}

	stdprinter.Debug("will cache index")
	err = e.Index.Write(f)

	if err != nil {
		return
	}

	stdprinter.Debug("did cache index")
	return
}

func (u Umwelt) Dir() string {
	return u.BasePath
}
