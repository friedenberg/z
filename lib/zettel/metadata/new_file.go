package metadata

import (
	"fmt"
	"path"

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
	if len(s) < 4 {
		err = xerrors.Errorf("string %s is too small to be a file tag", s)
		return
	}

	fd.Path = s[3:]

	return
}

func (fd NewFile) Ext() string {
	return path.Ext(fd.Path)
}

func (fd NewFile) Extension() string {
	return path.Ext(fd.Path)
}

func (fd NewFile) FilePath(_ string) string {
	//TODO-P1 should be relative?
	return fd.Path
}

func (fd NewFile) Tag() string {
	return fmt.Sprintf("%s-%s", NewFilePrefix, fd.Path)
}

func (f NewFile) SearchMatchTags() (expanded TagSet) {
	expanded = MakeTagSet()
	return
}

func (f NewFile) Match(_ string) bool {
	return false
}
