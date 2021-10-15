package printer

import (
	"encoding/json"

	"github.com/friedenberg/z/lib"
	"github.com/friedenberg/z/util/stdprinter"
	"github.com/pelletier/go-toml/v2"
)

type TomlToJson struct{}

func (p *TomlToJson) Begin() {}
func (p *TomlToJson) End()   {}

func (p *TomlToJson) PrintZettel(i int, z *lib.Zettel, errIn error) {
	if errIn != nil {
		stdprinter.Err(errIn)
		return
	}

	if _, ok := z.Metadata.TagSet().Get("k-toml"); !ok {
		return
	}

	var doc map[string]interface{}
	errOut := toml.Unmarshal([]byte(z.Body), &doc)

	if errOut != nil {
		stdprinter.Err(errOut)
		return
	}

	out := map[string]interface{}{
		"metadata": z.Metadata,
		"body":     doc,
	}

	b, errOut := json.Marshal(out)

	if errOut != nil {
		stdprinter.Err(errOut)
		return
	}

	stdprinter.Out(string(b))
}
