package lib

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path"

	"github.com/pelletier/go-toml/v2"
)

type ConfigTagForNewZettels struct {
	TagForNewZettels string `toml:"tag-for-new-zettels"`
}

type KastenImplementationType string

type KastenConfig struct {
	ConfigTagForNewZettels

	Implementation KastenImplementationType
	Options        map[string]string
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
	f, err := os.Open(p)

	if err != nil {
		return
	}

	defer f.Close()

	doc, err := ioutil.ReadAll(f)

	defer func() {
		if r := recover(); r != nil {
			c = Config{}
			err = fmt.Errorf("toml unmarshalling paniced: %q", r)
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

func (c Config) GetKasten() (ks []*Kasten, err error) {
	usr, err := user.Current()

	if err != nil {
		return
	}

	e := &Kasten{
		BasePath: path.Join(usr.HomeDir, "Zettelkasten"),
		Index:    MakeIndex(),
	}

	ks = append(ks, e)

	return
}
