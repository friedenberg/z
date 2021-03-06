package metadata

import (
	"fmt"
	"path"
	"strings"

	"github.com/friedenberg/z/util"
	"golang.org/x/xerrors"
)

func init() {
	registerTagPrefix(
		"f",
		func() (t ITag) { return &LocalFile{} },
	)
}

type LocalFile struct {
	FullSha string
	Id      string
	Ext     string
}

type RemoteFileDescriptor struct {
	LocalFile
}

func (fd *LocalFile) Set(s string) (err error) {
	if len(s) < 3 {
		err = xerrors.Errorf("string %s is too small to be a file tag", s)
		return
	}

	parts := strings.Split(s, "-")
	partCount := len(parts)

	if partCount != 2 {
		err = xerrors.Errorf("wrong number of tag parts: %s", partCount)
		return
	}

	fd.Id = util.BaseNameNoSuffix(parts[1])
	fd.Ext = util.ExtNoDot(parts[1])

	return
}

func (fd LocalFile) Tag() string {
	sb := &strings.Builder{}

	sb.WriteString("f-")

	sb.WriteString(fd.Id)

	if fd.Ext != "" {
		sb.WriteString(".")
		sb.WriteString(fd.Ext)
	}

	return sb.String()
}

func (f LocalFile) SearchMatchTags() (expanded TagSet) {
	expanded = MakeTagSet()

	if f.Ext != "" {
		expanded.Merge(Tag("e-" + f.Ext).SearchMatchTags())
	}

	return
}

func (fd LocalFile) FileName() (fn string) {
	fi := fd.Id

	if fd.Ext == "" {
		fn = fi
	} else {
		fn = fmt.Sprintf("%s.%s", fi, fd.Ext)
	}

	return
}

func (fd LocalFile) Extension() string {
	return fd.Ext
}

func (fd LocalFile) FilePath(basepath string) (fn string) {
	return path.Join(basepath, fd.FileName())
}
