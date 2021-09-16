package lib

import (
	"os/user"
	"path"
	"path/filepath"
	"time"

	"github.com/friedenberg/z/lib/kasten"
	"github.com/friedenberg/z/util"
)

var (
	FilesAndGitInstance *FilesAndGit
)

func init() {
	FilesAndGitInstance = &FilesAndGit{}
	kasten.Registry.Register("files-and-git", FilesAndGitInstance)
}

type FilesAndGit struct {
	BasePath        string
	GitEnabled      bool
	GitAnnexEnabled bool
}

func (f *FilesAndGit) getBoolOption(o map[string]interface{}, k string) bool {
	if s, ok := o[k]; ok {
		if sb, ok := s.(bool); ok {
			return sb
		}
	}

	//TODO
	//https://github.com/mitchellh/mapstructure
	return false
}

func (k *FilesAndGit) InitFromOptions(o map[string]interface{}) (err error) {
	usr, err := user.Current()

	if err != nil {
		return
	}

	k.BasePath = path.Join(usr.HomeDir, "Zettelkasten")

	k.GitEnabled = k.getBoolOption(o, "git-enabled")
	k.GitAnnexEnabled = k.getBoolOption(o, "git-annex-enabled")

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

func (k *FilesAndGit) NewId() (id Id, err error) {
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
			id = Id(currentTime.Unix())
			return
		}
	}

	return
}
