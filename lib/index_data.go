package lib

import (
	"bufio"
	"fmt"
	"io"
	"net/url"
	"os"
	"path"
	"sort"
	"strings"
	"time"

	"github.com/friedenberg/z/util"
	"gopkg.in/yaml.v2"
)

const (
	MetadataStartSequence = "---\n"
	MetadataEndSequence   = "...\n"
)

type Metadata []string

type ZettelIndexData struct {
	Date         string   `yaml:"-" json:"date,omitempty"`
	Description  string   `yaml:"description,omitempty" json:"description,omitempty"`
	Tags         []string `yaml:"tags,omitempty" json:"tags,omitempty"`
	ExpandedTags []string `yaml:"-" json:"expanded_tags,omitempty"`
	Url          string   `yaml:"url,omitempty" json:"url,omitempty"`
	File         string   `yaml:"file,omitempty" json:"file,omitempty"`
}

func (id ZettelIndexData) ToMetadata() (md Metadata) {
	if id.Description != "" {
		md = append(md, id.Description)
	}

	sort.Slice(id.Tags, func(i, j int) bool {
		return id.Tags[i] < id.Tags[j]
	})

	if len(id.Tags) > 0 {
		md = append(md, id.Tags...)
	}

	if id.File != "" {
		md = append(md, id.File)
	}

	if id.Url != "" {
		md = append(md, id.Url)
	}

	return
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

		if err == io.EOF {
			return nil
		}

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
	var md Metadata
	err = yaml.Unmarshal([]byte(z.Data.MetadataYaml), &md)

	if err != nil {
		err = fmt.Errorf("parse metadata: %w", err)
		return
	}

	for i, v := range md {
		if i == 0 {
			z.IndexData.Description = v
			continue
		}

		filePath := path.Join(z.Kasten.BasePath, v)

		if util.FileExists(filePath) {
			if z.IndexData.File != "" {
				err = fmt.Errorf(
					"zettel has more than one valid file: '%s' and '%s'",
					v,
					z.IndexData.File,
				)

				return
			}

			z.IndexData.File = v
			continue
		}

		url, e := url.Parse(v)

		if e == nil && url.Hostname() != "" {
			if z.IndexData.Url != "" {
				err = fmt.Errorf(
					"zettel has more than one valid url: '%s' and '%s'",
					v,
					z.IndexData.Url,
				)

				return
			}
			z.IndexData.Url = url.String()
			continue
		}

		z.IndexData.Tags = append(z.IndexData.Tags, v)
	}

	// if z.HasFile() {
	// 	var np string
	// 	np, err = z.Kasten.GetNormalizedPath(z.IndexData.File)

	// 	if err != nil {
	// 		return
	// 	}

	// 	z.IndexData.File = np
	// }

	var t time.Time

	t, err = TimeFromPath(z.Path)

	if err != nil {
		err = fmt.Errorf("parse metadata: %w", err)
		return
	}

	z.IndexData.Date = t.Format("2006-01-02")

	for _, t := range z.IndexData.Tags {
		z.IndexData.ExpandedTags = append(
			z.IndexData.ExpandedTags,
			util.ExpandTags(t)...,
		)
	}

	return
}
