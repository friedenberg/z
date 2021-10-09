package lib

import (
	"bufio"
	"bytes"

	"github.com/friedenberg/z/lib/zettel/metadata"
	"github.com/friedenberg/z/util/files_guard"
)

func (z *Zettel) ReadMetadataAndBody() (err error) {
	f, err := files_guard.Open(z.Path)
	defer files_guard.Close(f)

	if err != nil {
		return
	}

	r := bufio.NewReader(f)

	yamlString, err := metadata.ReadYAMLHeader(r)

	if err != nil {
		return
	}

	z.Data.MetadataYaml = yamlString

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
