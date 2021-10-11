package lib

import (
	"bufio"
	"os"

	"github.com/friedenberg/z/util/files_guard"
	"golang.org/x/xerrors"
)

func (z *Zettel) Write(onWriteFunc OnZettelWriteFunc) (err error) {
	if onWriteFunc != nil {
		defer onWriteFunc(z, err)
	}

	y, err := z.Note.Metadata.ToYAMLWithBoundary()

	if err != nil {
		return
	}

	f, err := files_guard.OpenFile(z.Path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	defer files_guard.Close(f)

	if err != nil {
		err = xerrors.Errorf("opening zettel file: %w", err)
		return
	}

	w := bufio.NewWriter(f)

	_, err = w.WriteString(y)

	if err != nil {
		return
	}

	if z.Body != "" {
		_, err = w.WriteString(z.Body)

		if err != nil {
			err = xerrors.Errorf("writing body: %w", err)
			return
		}
	}

	if err == nil {
		w.Flush()
	}

	return
}
