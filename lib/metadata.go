package lib

import (
	"bufio"
	"net/url"
	"path"
	"regexp"
	"sort"

	"github.com/friedenberg/z/lib/zettel/metadata"
	"github.com/friedenberg/z/util"
	"github.com/friedenberg/z/util/files_guard"
	"golang.org/x/xerrors"
	"gopkg.in/yaml.v2"
)

var (
	RegexTag *regexp.Regexp
)

func init() {
	//TODO decide on regex
	RegexTag = regexp.MustCompile(`^\w{1,2}-[^\s]+$`)
}

type MetadataList []string

type Metadata struct {
	Description  string   `yaml:"description,omitempty" json:"description,omitempty"`
	Tags         []string `yaml:"tags,omitempty" json:"tags,omitempty"`
	ExpandedTags []string `yaml:"-" json:"expanded_tags,omitempty"`
	Url          string   `yaml:"url,omitempty" json:"url,omitempty"`
	File         string   `yaml:"file,omitempty" json:"file,omitempty"`
}

func (id Metadata) ToMetadata() (md MetadataList) {
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

func (z *Zettel) ReadMetadata() (err error) {
	f, err := files_guard.Open(z.Path)
	defer files_guard.Close(f)

	if err != nil {
		return
	}

	r := bufio.NewReader(f)

	yamlString, err := metadata.ReadYAMLHeader(r)

	if err != nil {
		return
	}

	z.Data.MetadataYaml = yamlString

	return
}

func (z *Zettel) ParseMetadata() (err error) {
	var md MetadataList
	err = yaml.Unmarshal([]byte(z.Data.MetadataYaml), &md)

	if err != nil {
		err = xerrors.Errorf("parse metadata: %w", err)
		return
	}

	return z.FromMetadata(md)
}

func (z *Zettel) FromMetadata(md MetadataList) (err error) {
	for i, v := range md {
		if i == 0 && !RegexTag.MatchString(v) {
			z.Metadata.Description = v
			continue
		}

		filePath := path.Join(z.Kasten.Local.BasePath, v)

		if util.FileExists(filePath) {
			if z.Metadata.File != "" {
				err = xerrors.Errorf(
					"zettel has more than one valid file: '%s' and '%s'",
					v,
					z.Metadata.File,
				)

				return
			}

			z.Metadata.File = v
			continue
		}

		url, e := url.Parse(v)

		if e == nil && url.Hostname() != "" {
			if z.Metadata.Url != "" {
				err = xerrors.Errorf(
					"zettel has more than one valid url: '%s' and '%s'",
					v,
					z.Metadata.Url,
				)

				return
			}
			z.Metadata.Url = url.String()
			continue
		}

		z.Metadata.Tags = append(z.Metadata.Tags, v)
	}

	_, err = TimeFromPath(z.Path)

	if err != nil {
		err = xerrors.Errorf("parse metadata: %w", err)
		return
	}

	// z.Metadata.Date = t.Format("2006-01-02")

	for _, t := range z.Metadata.Tags {
		z.Metadata.ExpandedTags = append(
			z.Metadata.ExpandedTags,
			metadata.ExpandTags(t)...,
		)
	}

	if z.HasFile() {
		z.Metadata.ExpandedTags = append(z.Metadata.ExpandedTags, "h-f")
	}

	if z.HasUrl() {
		z.Metadata.ExpandedTags = append(z.Metadata.ExpandedTags, "h-u")
	}

	return
}
