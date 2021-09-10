package lib

import (
	"os"
	"os/user"
	"path"

	"github.com/friedenberg/z/lib/kasten"
)

type Umwelt struct {
	DefaultKasten kasten.Implementation
	Kasten        map[string]kasten.Implementation
	Index         Index
	BasePath      string
	Config        Config
}

func MakeUmwelt(c Config) (k Umwelt, err error) {
	k.Config = c

	usr, err := user.Current()

	if err != nil {
		return
	}

	k.BasePath = path.Join(usr.HomeDir, "Zettelkasten")
	k.Index = MakeIndex()

	return
}

func (u Umwelt) FilesAndGit() *FilesAndGit {
	return u.DefaultKasten.(*FilesAndGit)
}
