package lib

import (
	"os"
	"os/user"
	"path"

	"github.com/friedenberg/z/lib/kasten"
	"github.com/friedenberg/z/util/files_guard"
)

type Umwelt struct {
	Kasten       Kasten
	RemoteKasten map[string]kasten.RemoteImplementation
	Index        Index
	BasePath     string
	Config       Config
}

func MakeUmwelt(c Config) (k Umwelt, err error) {
	k.Config = c

	usr, err := user.Current()

	if err != nil {
		return
	}

	k.BasePath = path.Join(usr.HomeDir, "Zettelkasten")
	k.Index = MakeIndex()

	if c.UseIndexCache {
		//find all caches
		//remove any older-version caches
		//try to load from correct version cache
		//or exit if it doesn't exist
		err = k.LoadIndexFromCache()

		if err != nil && !os.IsNotExist(err) {
			return
		}
	}

	return
}

func (u Umwelt) FilesAndGit() *FilesAndGit {
	return u.Kasten.LocalImplementation.(*FilesAndGit)
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
