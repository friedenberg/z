package lib

import (
	"bufio"
	"bytes"

	"github.com/friedenberg/z/util/files_guard"
)

func (zettel *KastenZettel) ReadMetadataAndBody() (err error) {
	f, err := files_guard.Open(zettel.Path)
	defer files_guard.Close(f)

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

func (z *KastenZettel) readBodyFromReader(r *bufio.Reader) (err error) {
	body := &bytes.Buffer{}
	_, err = r.WriteTo(body)

	if err != nil {
		return
	}

	z.Body = string(body.Bytes())

	return
}
