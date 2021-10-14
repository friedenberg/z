package lib

import (
	"io/ioutil"
	"os"
	"os/user"
	"path"

	"github.com/friedenberg/z/lib/kasten"
	"github.com/friedenberg/z/util/files_guard"
	"github.com/pelletier/go-toml/v2"
	"golang.org/x/xerrors"
)

type ConfigTagForNewZettels struct {
	TagForNewZettels string `toml:"tag-for-new-zettels"`
}

type KastenRemoteConfig struct {
	ConfigTagForNewZettels

	Implementation string `toml:"type"`
	Options        map[string]interface{}
}

type TagConfig struct {
	kasten          string
	AutoTags        []string `toml:"auto-tags"`
	AddToNewZettels bool     `toml:"add-to-new-zettels"`
}

type Config struct {
	ConfigTagForNewZettels

	Path           string `toml:"path"`
	GitEnabled     bool   `toml:"git-enabled"`
	GitSignCommits bool   `toml:"git-sign-commits"`
	Remotes        map[string]KastenRemoteConfig
	Tags           map[string]TagConfig
}

func DefaultConfigPath() (p string, err error) {
	usr, err := user.Current()

	if err != nil {
		return
	}

	p = path.Join(
		usr.HomeDir,
		".config",
		"zettelkasten",
		"config.toml",
	)

	return
}

func LoadConfig(p string) (c Config, err error) {
	f, err := files_guard.Open(p)
	defer files_guard.Close(f)

	if err != nil {
		return
	}

	doc, err := ioutil.ReadAll(f)

	defer func() {
		if r := recover(); r != nil {
			c = Config{}
			err = xerrors.Errorf("toml unmarshalling paniced: %q", r)
		}
	}()

	err = toml.Unmarshal([]byte(doc), &c)

	if err != nil {
		return
	}

	return
}

func DefaultConfig() (c Config, err error) {
	//TODO-P2
	return
}

func LoadDefaultConfig() (c Config, err error) {
	p, err := DefaultConfigPath()

	if err != nil {
		return
	}

	c, err = LoadConfig(p)

	if os.IsNotExist(err) {
		c, err = DefaultConfig()
		return
	} else if err != nil {
		return
	}

	return
}

func (c Config) Umwelt() (u Umwelt, err error) {
	u, err = MakeUmwelt(c)

	if err != nil {
		return
	}

	wd, err := os.Getwd()

	if err != nil {
		return
	}

	fs := &FileStore{
		//TODO-P2 use cwd or config if available
		basePath: wd,
	}

	if c.GitEnabled {
		u.Kasten.Local = &GitStore{
			FileStore:   *fs,
			SignCommits: c.GitSignCommits,
		}
	} else {
		u.Kasten.Local = fs
	}

	u.Kasten.Remotes = make(map[string]kasten.RemoteImplementation)

	for n, kc := range c.Remotes {
		if i, ok := kasten.GetRemote(kc.Implementation); ok {
			i.InitFromOptions(kc.Options)
			u.Kasten.Remotes[n] = i
		} else {
			err = xerrors.Errorf("missing implementation for kasten from config: '%s'", n)
			return
		}
	}

	for t, tc := range c.Tags {
		if tc.AddToNewZettels {
			u.TagsForNewZettels = append(u.TagsForNewZettels, t)
		}
	}

	return
}
