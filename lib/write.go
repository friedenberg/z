package lib

import (
	"bufio"
	"os"

	"github.com/friedenberg/z/lib/zettel/metadata"
	"github.com/friedenberg/z/util/files_guard"
	"golang.org/x/xerrors"
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
		err = xerrors.Errorf("writing zettel: %w", err)
		return
	}

	f, err := files_guard.OpenFile(z.Path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	defer files_guard.Close(f)

	if err != nil {
		err = xerrors.Errorf("opening zettel file: %w", err)
		return
	}

	w := bufio.NewWriter(f)

	_, err = w.WriteString(metadata.MetadataStartSequence)

	if err != nil {
		err = xerrors.Errorf("writing metadata start sequence: %w", err)
		return
	}

	_, err = w.WriteString(z.Data.MetadataYaml)

	if err != nil {
		err = xerrors.Errorf("writing metadata yaml: %w", err)
		return
	}

	_, err = w.WriteString(metadata.MetadataEndSequence)

	if err != nil {
		err = xerrors.Errorf("writing metadata end sequence: %w", err)
		return
	}

	if z.Body != "" {
		_, err = w.WriteString(z.Body)

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
