package lib

import (
	"bufio"
	"os"

	"gopkg.in/yaml.v2"
)

func (z *Zettel) Write() (err error) {
	var y []byte
	y, err = yaml.Marshal(z.Metadata)

	if err != nil {
		return
	}

	z.MetadataYaml = string(y)

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

	_, err = w.WriteString(z.MetadataYaml)

	if err != nil {
		return
	}

	_, err = w.WriteString(MetadataEndSequence)

	if err != nil {
		return
	}

	if z.Body == "" {
		return
	}

	_, err = w.WriteString("\n")

	if err != nil {
		return
	}

	_, err = w.WriteString(z.Body)

	return
}
