package metadata

import (
	"path"
	"strings"

	"golang.org/x/xerrors"
)

const NewFilePrefix = "nf"

func init() {
	registerTagPrefix(
		NewFilePrefix,
		func() (t ITag) { return &NewFile{} },
	)
}

type NewFile struct {
	Path string
}

func (fd *NewFile) Set(s string) (err error) {
	if len(s) < 3 {
		err = xerrors.Errorf("string %s is too small to be a file tag", s)
		return
	}

	parts := strings.Split(s, "-")
	partCount := len(parts)

	if partCount > 3 || partCount < 2 {
		err = xerrors.Errorf("wrong number of tag parts: %s", partCount)
		return
	}

	fd.Path = parts[1]

	return
}

func (fd NewFile) Ext() string {
	return path.Ext(fd.Path)
}

func (fd NewFile) Tag() string {
	return ""
}

func (f NewFile) SearchMatchTags() (expanded TagSet) {
	expanded = MakeTagSet()
	return
}
