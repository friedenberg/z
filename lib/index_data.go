package lib

import (
	"bufio"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

const (
	MetadataStartSequence = "---\n"
	MetadataEndSequence   = "...\n"
)

type ZettelIndexData struct {
	Date        string   `yaml:"date,omitempty" json:"date,omitempty"`
	Description string   `yaml:"description,omitempty" json:"description,omitempty"`
	Areas       []string `yaml:"areas,omitempty" json:"areas,omitempty"`
	Projects    []string `yaml:"projects,omitempty" json:"projects,omitempty"`
	Tags        []string `yaml:"tags,omitempty" json:"tags,omitempty"`
	Url         string   `yaml:"url,omitempty" json:"url,omitempty"`
	File        string   `yaml:"file,omitempty" json:"file,omitempty"`
}

func (zettel *Zettel) ReadMetadata() (err error) {
	f, err := os.Open(zettel.Path)
	defer f.Close()

	if err != nil {
		return
	}

	r := bufio.NewReader(f)

	return zettel.readMetadataFromReader(r)
}

func (z *Zettel) readMetadataFromReader(r *bufio.Reader) (err error) {
	sb := strings.Builder{}
	within := false

	for {
		some_string, err := r.ReadString('\n')

		if err != nil {
			return err
		}

		if !within && some_string == MetadataStartSequence {
			within = true
		} else if within && some_string == MetadataEndSequence {
			break
		} else if within {
			sb.WriteString(some_string)
		}
	}

	z.Data.MetadataYaml = sb.String()

	return
}

func (z *Zettel) ParseMetadata() (err error) {
	err = yaml.Unmarshal([]byte(z.Data.MetadataYaml), &z.IndexData)

	if z.HasFile() {
		var np string
		np, err = z.Env.GetNormalizedPath(z.IndexData.File)

		if err != nil {
			return
		}

		z.IndexData.File = np
	}

	return
}
