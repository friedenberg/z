package lib

import (
	"os"
	"path"

	"github.com/friedenberg/z/util/files_guard"
)

type Umwelt struct {
	Kasten
	Index             Index
	BasePath          string
	Config            Config
	TagsForNewZettels []string
	Transaction
}

func MakeUmwelt(c Config) (k Umwelt, err error) {
	k.Config = c

	wd, err := os.Getwd()

	if err != nil {
		return
	}

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

func (u Umwelt) Store() Store {
	return u.Kasten.Local
}

func (e Umwelt) GetIndexPath() string {
	return path.Join(e.BasePath, ".zettel-cache")
}

func (u Umwelt) GetAll() (ids []string) {
	ids = make([]string, 0, len(u.Index.Zettels))

	for id, _ := range u.Index.Zettels {
		ids = append(ids, id.String())
	}

	return
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

	err = e.Index.Write(f)

	if err != nil {
		return
	}

	return
}
