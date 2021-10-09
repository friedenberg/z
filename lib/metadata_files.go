package lib

import (
	"path"
	"strconv"

	"github.com/friedenberg/z/lib/zettel/metadata"
	"golang.org/x/xerrors"
)

func (z *Zettel) LocalFile() (fd metadata.File, ok bool) {
	ok = z.HasFile()

	if !ok {
		return
	}

	fd = metadata.File{
		Id:  strconv.FormatInt(z.Id, 10),
		Ext: path.Ext(z.FilePath()),
	}

	return
}

func (z *Zettel) RemoteFiles() (fds []metadata.File) {
	return
}

func (z *Zettel) AddFileDescripter(fd metadata.File) (err error) {
	z.Metadata.Tags = append(z.Metadata.Tags, fd.Tag())
	return
}

func (z *Zettel) RemoveFileDescripter(fd metadata.File) (err error) {
	found := -1
	tags := z.Metadata.Tags
	tag := fd.Tag()

	for i, t := range tags {
		if t == tag {
			found = i
			break
		}
	}

	if found == -1 {
		err = xerrors.Errorf("tag not found: %s", tag)
		return
	}

	tags[found] = tags[len(tags)-1]
	z.Metadata.Tags = tags[:len(tags)-1]

	return
}
