package lib

import (
	"bufio"
	"os"

	"gopkg.in/yaml.v2"
)

type OnZettelWriteFunc func(*Zettel, error) error

func (z *Zettel) Write(onWriteFunc OnZettelWriteFunc) (err error) {
	if onWriteFunc != nil {
		defer onWriteFunc(z, err)
	}

	var y []byte
	y, err = yaml.Marshal(z.IndexData)

	if err != nil {
		return
	}

	z.Data.MetadataYaml = string(y)

	//TODO
	f, err := os.OpenFile(z.Path, os.O_RDWR|os.O_CREATE, 0755)
	defer f.Close()

	if err != nil {
		return
	}

	w := bufio.NewWriter(f)
	defer w.Flush()

	_, err = w.WriteString(MetadataStartSequence)

	if err != nil {
		return
	}

	_, err = w.WriteString(z.Data.MetadataYaml)

	if err != nil {
		return
	}

	_, err = w.WriteString(MetadataEndSequence)

	if err != nil {
		return
	}

	if z.Data.Body == "" {
		return
	}

	_, err = w.WriteString("\n")

	if err != nil {
		return
	}

	_, err = w.WriteString(z.Data.Body)

	return
}
