package lib

import (
	"bufio"
	"fmt"
	"os"

	"github.com/friedenberg/z/util"
	"gopkg.in/yaml.v2"
)

func (z *Zettel) generateMetadataYaml() (err error) {
	var y []byte
	y, err = yaml.Marshal(z.Metadata.ToMetadata())

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

	err = z.generateMetadataYaml()

	if err != nil {
		err = fmt.Errorf("writing zettel: %w", err)
		return
	}

	//TODO
	util.OpenFilesGuardInstance.Lock()
	f, err := os.OpenFile(z.Path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	defer f.Close()
	defer util.OpenFilesGuardInstance.Unlock()

	if err != nil {
		err = fmt.Errorf("opening zettel file: %w", err)
		return
	}

	w := bufio.NewWriter(f)

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

	if z.Data.Body != "" {
		_, err = w.WriteString(z.Data.Body)

		if err != nil {
			err = fmt.Errorf("writing body: %w", err)
			return
		}
	}

	if err == nil {
		w.Flush()
	}

	return
}
