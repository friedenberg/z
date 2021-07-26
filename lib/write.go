package lib

import (
	"bufio"
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

type OnZettelWriteFunc func(*Zettel, error) error

func (z *Zettel) GenerateMetadataYaml() (err error) {
	var y []byte
	y, err = yaml.Marshal(z.IndexData)

	if err != nil {
		return
	}

	z.Data.MetadataYaml = string(y)

	return
}

func (z *Zettel) Write(onWriteFunc OnZettelWriteFunc) (err error) {
	if onWriteFunc != nil {
		defer onWriteFunc(z, err)
	}

	err = z.GenerateMetadataYaml()

	if err != nil {
		err = fmt.Errorf("writing zettel: %w", err)
		return
	}

	//TODO
	f, err := os.OpenFile(z.Path, os.O_RDWR|os.O_CREATE, 0755)
	defer f.Close()

	if err != nil {
		err = fmt.Errorf("opening zettel file: %w", err)
		return
	}

	w := bufio.NewWriter(f)
	defer w.Flush()

	_, err = w.WriteString(MetadataStartSequence)

	if err != nil {
		err = fmt.Errorf("writing metadata start sequence: %w", err)
		return
	}

	_, err = w.WriteString(z.Data.MetadataYaml)

	if err != nil {
		err = fmt.Errorf("writing metadata yaml: %w", err)
		return
	}

	_, err = w.WriteString(MetadataEndSequence)

	if err != nil {
		err = fmt.Errorf("writing metadata end sequence: %w", err)
		return
	}

	if z.Data.Body == "" {
		return
	}

	_, err = w.WriteString(z.Data.Body)

	if err != nil {
		err = fmt.Errorf("writing body: %w", err)
		return
	}

	return
}
