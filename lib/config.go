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

type KastenConfig struct {
	ConfigTagForNewZettels

	Implementation string `toml:"type"`
	Options        map[string]interface{}
}

type TagConfig struct {
	kasten   string
	autoTags []string `toml:"auto-tags"`
}

type Config struct {
	ConfigTagForNewZettels

	Tags          map[string]TagConfig
	Kasten        map[string]KastenConfig
	DefaultKasten string `toml:"default-kasten"`
	UseIndexCache bool   `toml:"use-index-cache"`
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
	//TODO
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

func (c Config) Umwelt() (e Umwelt, err error) {
	e, err = MakeUmwelt(c)

	if err != nil {
		return
	}

	e.Kasten = make(map[string]kasten.Implementation)

	for n, kc := range c.Kasten {
		if i, ok := kasten.Registry.Get(kc.Implementation); ok {
			i.InitFromOptions(kc.Options)
			e.Kasten[n] = i
			//TODO
			e.DefaultKasten = i
		} else {
			err = xerrors.Errorf("missing implementation for kasten from config: '%s'", n)
			return
		}
	}

	if c.DefaultKasten != "" {
		if i, ok := e.Kasten[c.DefaultKasten]; ok {
			e.DefaultKasten = i
		} else {
			err = xerrors.Errorf(
				"no kasten matching name '%s' for default",
				c.DefaultKasten,
			)

			return
		}
	}

	return
}
