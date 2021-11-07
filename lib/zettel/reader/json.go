package reader

import (
	"encoding/json"

	"github.com/friedenberg/z/lib"
	"github.com/friedenberg/z/lib/zettel"
	"github.com/friedenberg/z/util/stdprinter"
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

	//TODO-P3 try to read ID or assign ID
	err = json.Unmarshal(b, &z)
	stdprinter.Debugf("%#v", z)

	if err != nil {
		err = xerrors.Errorf("unable to read zettel from json: %w", err)
		return
	}

	return
}
