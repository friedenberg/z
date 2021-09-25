package lib

import (
	"bufio"
	"io"
	"net/url"
	"path"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/friedenberg/z/lib/zettel"
	"github.com/friedenberg/z/util"
	"github.com/friedenberg/z/util/files_guard"
	"golang.org/x/xerrors"
	"gopkg.in/yaml.v2"
)

const (
	MetadataStartSequence = "---\n"
	MetadataEndSequence   = "...\n"
)

var (
	RegexTag *regexp.Regexp
)

func init() {
	RegexTag = regexp.MustCompile(`^\w{1,2}-[^\s]+$`)
}

type MetadataList []string

type Metadata struct {
	Date         string   `yaml:"-" json:"date,omitempty"`
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

func (zettel *KastenZettel) ReadMetadata() (err error) {
	f, err := files_guard.Open(zettel.Path)
	defer files_guard.Close(f)

	if err != nil {
		return
	}

	r := bufio.NewReader(f)

	return zettel.readMetadataFromReader(r)
}

func (z *KastenZettel) readMetadataFromReader(r *bufio.Reader) (err error) {
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

func (z *KastenZettel) ParseMetadata() (err error) {
	var md MetadataList
	err = yaml.Unmarshal([]byte(z.Data.MetadataYaml), &md)

	if err != nil {
		err = xerrors.Errorf("parse metadata: %w", err)
		return
	}

	return z.FromMetadata(md)
}

func (z *KastenZettel) FromMetadata(md MetadataList) (err error) {
	for i, v := range md {
		if i == 0 && !RegexTag.MatchString(v) {
			z.Metadata.Description = v
			continue
		}

		filePath := path.Join(z.Umwelt.FilesAndGit().BasePath, v)

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

	var t time.Time

	t, err = TimeFromPath(z.Path)

	if err != nil {
		err = xerrors.Errorf("parse metadata: %w", err)
		return
	}

	z.Metadata.Date = t.Format("2006-01-02")

	for _, t := range z.Metadata.Tags {
		z.Metadata.ExpandedTags = append(
			z.Metadata.ExpandedTags,
			util.ExpandTags(t)...,
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

func (z *KastenZettel) FileDescriptor() (fd *zettel.FileDescriptor) {
	if !z.HasFile() {
		return
	}
	fd = &zettel.FileDescriptor{
		ZettelId: zettel.Id(z.Id),
		Ext:      path.Ext(z.FilePath()),
	}

	return
}
