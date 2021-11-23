package reader

import (
	"encoding/json"

	"github.com/friedenberg/z/lib"
	"github.com/friedenberg/z/lib/zettel"
	"golang.org/x/xerrors"
)

type Json struct {
}

func (p *Json) ReadZettel(u *lib.Umwelt, i int, b []byte) (z *zettel.Zettel, err error) {
	z, err = readerNew(u, i, "")

	if err != nil {
		err = xerrors.Errorf("unable to make new zettel: %w", err)
		return
	}

	err = json.Unmarshal(b, &z)

	if err != nil {
		err = xerrors.Errorf("unable to read zettel from json: %w", err)
		return
	}

	return
}
