package zettel

import (
	"bufio"
	"io"
	"os"

	"github.com/friedenberg/z/util/files_guard"
	"golang.org/x/xerrors"
)

type OnZettelWriteFunc func(*Zettel, error) error

func (z *Zettel) Write(onWriteFunc OnZettelWriteFunc) (err error) {
	if onWriteFunc != nil {
		defer onWriteFunc(z, err)
	}

	f, err := files_guard.OpenFile(z.Path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	defer files_guard.Close(f)

	if err != nil {
		err = xerrors.Errorf("opening zettel file: %w", err)
		return
	}

	w := bufio.NewWriter(f)

	err = z.WriteTo(w)

	if err != nil {
		return
	}

	return
}

func (z *Zettel) WriteTo(w1 io.Writer) (err error) {
	y, err := z.Note.Metadata.ToYAMLWithBoundary()
	w := bufio.NewWriter(w1)

	if err != nil {
		return
	}

	if err != nil {
		err = xerrors.Errorf("opening zettel file: %w", err)
		return
	}

	_, err = w.WriteString(y)

	if err != nil {
		return
	}

	if z.Body != "" {
		_, err = w.WriteString("\n")

		if err != nil {
			err = xerrors.Errorf("writing body: %w", err)
			return
		}

		_, err = w.WriteString(z.Body)

		if err != nil {
			err = xerrors.Errorf("writing body: %w", err)
			return
		}

		_, err = w.WriteString("\n")

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
