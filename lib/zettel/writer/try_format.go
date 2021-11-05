package writer

import (
	"encoding/json"
	"io"

	"github.com/friedenberg/z/lib/zettel"
	"github.com/friedenberg/z/util/stdprinter"
)

type TryFormat struct{}

func (p *TryFormat) WriteZettel(w io.Writer, i int, z *zettel.Zettel) {
	var err error
	var doc interface{}
	if _, ok := z.Metadata.TagSet().Get("k-toml"); ok {
		doc, err = TomlZettelBody(z)

		if err != nil {
			//TODO
			return
		}
	} else {
		doc = z.Body
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

	_, err = w.Write(b)
	stdprinter.PanicIfError(err)
}
