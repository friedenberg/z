package lib

import (
	"os/user"
	"path"
	"path/filepath"

	"github.com/friedenberg/z/lib/kasten"
)

var (
	FilesAndGitInstance *FilesAndGit
)

func init() {
	FilesAndGitInstance = &FilesAndGit{}
	kasten.Registry.Register("files-and-git", FilesAndGitInstance)
}

type FilesAndGit struct {
	BasePath string
	Index    Index
}

func (k *FilesAndGit) InitFromOptions(o map[string]string) (err error) {
	usr, err := user.Current()

	if err != nil {
		return
	}

	k.BasePath = path.Join(usr.HomeDir, "Zettelkasten")
	k.Index = MakeIndex()

	return
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
