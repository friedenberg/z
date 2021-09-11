package lib

import (
	"bufio"
	"bytes"

	"github.com/friedenberg/z/util"
)

func (zettel *Zettel) ReadMetadataAndBody() (err error) {
	f, err := util.OpenFilesGuardInstance.Open(zettel.Path)
	defer util.OpenFilesGuardInstance.Close(f)

	if err != nil {
		return
	}

	r := bufio.NewReader(f)

	err = zettel.readMetadataFromReader(r)

	if err != nil {
		return
	}

	err = zettel.readBodyFromReader(r)
	return
}

func (z *Zettel) readBodyFromReader(r *bufio.Reader) (err error) {
	body := &bytes.Buffer{}
	_, err = r.WriteTo(body)

	if err != nil {
		return
	}

	z.Body = string(body.Bytes())

	return
}
