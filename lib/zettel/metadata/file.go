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
		func() (t ITag) { return &File{} },
	)
}

const IdTruncationLength = 7

type File struct {
	KastenName string
	Id         string
	Ext        string
}

type RemoteFileDescriptor struct {
	File
}

func (fd *File) Set(s string) (err error) {
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

	fd.Id = util.BaseNameNoSuffix(parts[1])
	fd.Ext = util.ExtNoDot(parts[1])

	if partCount == 3 {
		fd.KastenName = parts[2]
	}

	return
}

func (fd File) TruncatedId() (i string) {
	i = fd.Id

	if len(i) > IdTruncationLength {
		i = i[0:IdTruncationLength]
	}

	return
}

func (fd File) Tag() string {
	sb := &strings.Builder{}

	sb.WriteString("f-")

	sb.WriteString(fd.TruncatedId())

	if fd.Ext != "" {
		sb.WriteString(".")
		sb.WriteString(fd.Ext)
	}

	if fd.KastenName != "" {
		sb.WriteString("-")
		sb.WriteString(fd.KastenName)
	}

	return sb.String()
}

func (f File) SearchMatchTags() (expanded TagSet) {
	expanded = MakeTagSet()

	if f.KastenName != "" {
		expanded.Merge(Tag("r-" + f.KastenName).SearchMatchTags())
	}

	if f.Ext != "" {
		expanded.Merge(Tag("e-" + f.Ext).SearchMatchTags())
	}

	return
}

func (fd File) FileName() (fn string) {
	fi := fd.TruncatedId()

	if fd.Ext == "" {
		fn = fi
	} else {
		fn = fmt.Sprintf("%s.%s", fi, fd.Ext)
	}

	return
}

func (fd File) FilePath(basepath string) (fn string) {
	return path.Join(basepath, fd.FileName())
}
