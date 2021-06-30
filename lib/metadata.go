package lib

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

const (
	MetadataStartSequence = "---\n"
	MetadataEndSequence   = "...\n"
)

type ZettelMetadata struct {
	Date        string   `yaml:"date,omitempty" json:"date,omitempty"`
	Kind        string   `yaml:"kind,omitempty" json:"kind,omitempty"`
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

	z.MetadataYaml = sb.String()

	return
}

func (zettel *Zettel) ParseMetadata() (err error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
		}
	}()

	// zettel.metadata.Description = zettel.metadataYaml
	err = yaml.Unmarshal([]byte(zettel.MetadataYaml), &zettel.Metadata)

	return
}
