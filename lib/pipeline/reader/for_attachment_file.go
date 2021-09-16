package reader

import (
	"os/exec"
	"strings"

	"github.com/friedenberg/z/lib"
	"github.com/friedenberg/z/lib/zettel/metadata"
	"github.com/friedenberg/z/util"
	"golang.org/x/xerrors"
)

func ForAttachmentFile(shouldCopy bool) reader {
	return Make(
		func(u lib.Umwelt, i int, b []byte) (*lib.Zettel, error) {
			return newOrFoundForFile(u, i, string(b), shouldCopy)
		},
	)
}

func newOrFoundForFile(u lib.Umwelt, i int, file string, shouldCopy bool) (z *lib.Zettel, err error) {
	sum, err := util.Sha256HashForFile(file)

	if err != nil {
		return
	}

	ids, ok := u.Index.Files.Get(sum, u.Index)

	if ok && ids.Len() > 1 {
		err = xerrors.Errorf("multiple zettels ('%q') with file: '%s'", ids, sum)
		return
	} else if ok && ids.Len() == 1 {
		z, err = FromIndex(u, 0, ids.Slice()[0].String())
		return
	}

	z, err = readerNew(u, i, file)

	if err != nil {
		return
	}

	//TODO-P0 check for checksum file name collisions
	n := sum[0:7]

	fd := metadata.File{
		Id:  n,
		Ext: strings.ReplaceAll(util.ExtNoDot(file), "-", ""),
	}

	if shouldCopy {
		cmd := exec.Command("cp", "-R", file, fd.FilePath(u.BasePath))
		msg, err := cmd.CombinedOutput()

		if err != nil {
			err = xerrors.Errorf("%w: %s", err, msg)
		}
	} else {
		cmd := exec.Command("mv", file, fd.FilePath(u.BasePath))
		msg, err := cmd.CombinedOutput()

		if err != nil {
			err = xerrors.Errorf("%w: %s", err, msg)
		}
	}

	if err != nil {
		return
	}

	z.Note.Metadata.AddFile(fd)

	return
}
