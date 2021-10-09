package metadata

import (
	"strings"

	"github.com/friedenberg/z/util"
	"golang.org/x/xerrors"
)

type File struct {
	KastenName string
	Id         string
	Ext        string
}

type RemoteFileDescriptor struct {
	File
}

func (fd *File) Set(s string) (err error) {
	parts := strings.Split(s, "-")
	partCount := len(parts)

	switch partCount {

	}

	if partCount > 2 || partCount < 1 {
		err = xerrors.Errorf("wrong number of tag parts: %s", partCount)
		return
	}

	fd.Id = util.BaseNameNoSuffix(parts[0])
	fd.Ext = util.ExtNoDot(parts[0])

	if partCount == 2 {
		fd.KastenName = parts[1]
	}

	return
}

func (fd File) Tag() string {
	sb := &strings.Builder{}

	sb.WriteString(fd.Id)

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

func (fd File) FileName() (fn string) {
	fi := fd.Id

	if fd.Ext == "" {
		fn = fi
	} else {
		fn = fi + fd.Ext
	}

	return
}
