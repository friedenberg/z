package metadata

import (
	"golang.org/x/xerrors"
	"gopkg.in/yaml.v2"
)

type Metadata struct {
	Description string
	Tags        []ITag
}

func (m *Metadata) Set(s string) (err error) {
	var tags []string
	err = yaml.Unmarshal([]byte(s), &tags)

	if err != nil {
		err = xerrors.Errorf("parse metadata: %w", err)
		return
	}

	// for i, v := range md {
	// 	if i == 0 && !RegexTag.MatchString(v) {
	// 		z.Metadata.Description = v
	// 		continue
	// 	}

	// 	filePath := path.Join(z.Kasten.Local.BasePath, v)

	// 	if util.FileExists(filePath) {
	// 		if z.Metadata.File != "" {
	// 			err = xerrors.Errorf(
	// 				"zettel has more than one valid file: '%s' and '%s'",
	// 				v,
	// 				z.Metadata.File,
	// 			)

	// 			return
	// 		}

	// 		z.Metadata.File = v
	// 		continue
	// 	}

	// 	url, e := url.Parse(v)

	// 	if e == nil && url.Hostname() != "" {
	// 		if z.Metadata.Url != "" {
	// 			err = xerrors.Errorf(
	// 				"zettel has more than one valid url: '%s' and '%s'",
	// 				v,
	// 				z.Metadata.Url,
	// 			)

	// 			return
	// 		}
	// 		z.Metadata.Url = url.String()
	// 		continue
	// 	}

	// 	z.Metadata.Tags = append(z.Metadata.Tags, v)
	// }

	// _, err = TimeFromPath(z.Path)

	// if err != nil {
	// 	err = xerrors.Errorf("parse metadata: %w", err)
	// 	return
	// }

	// // z.Metadata.Date = t.Format("2006-01-02")

	// for _, t := range z.Metadata.Tags {
	// 	z.Metadata.ExpandedTags = append(
	// 		z.Metadata.ExpandedTags,
	// 		metadata.ExpandTags(t)...,
	// 	)
	// }

	// if z.HasFile() {
	// 	z.Metadata.ExpandedTags = append(z.Metadata.ExpandedTags, "h-f")
	// }

	// if z.HasUrl() {
	// 	z.Metadata.ExpandedTags = append(z.Metadata.ExpandedTags, "h-u")
	// }

	return
}

func (m Metadata) Strings() (r []string) {
	return
}

func (m Metadata) WithExpandedTags() (m2 Metadata) {
	tags := make([]ITag, len(m.Tags))

	for _, t := range m.Tags {
		tags = append(tags, t)
		tags = append(tags, ExpandTagsAsTags(t.Tag())...)
	}

	m2.Description = m.Description
	m2.Tags = tags

	return
}
