package zettel

import "strings"

type FileDescriptor struct {
	KastenName string
	ZettelId   Id
	Ext        string
}

type RemoteFileDescriptor struct {
	FileDescriptor
}

func (fd *FileDescriptor) Set(s string) (err error) {
	return
}

func (fd FileDescriptor) Tag() string {
	sb := &strings.Builder{}

	sb.WriteString(fd.ZettelId.String())

	if fd.Ext != "" {
		sb.WriteString(fd.Ext)
	}

	if fd.KastenName != "" {
		sb.WriteString("-")
		sb.WriteString(fd.KastenName)
	}

	return sb.String()
}

func (fd FileDescriptor) FileName() (fn string) {
	fi := fd.ZettelId.String()

	if fd.Ext == "" {
		fn = fi
	} else {
		fn = fi + fd.Ext
	}

	return
}
