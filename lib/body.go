package lib

import (
	"bufio"
	"os"
)

func (zettel *Zettel) ReadMetadataAndBody() (err error) {
	f, err := os.Open(zettel.Path)
	defer f.Close()

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
	var body []byte
	_, err = r.Read(body)

	if err != nil {
		return
	}

	z.Body = string(body)
	return
}
