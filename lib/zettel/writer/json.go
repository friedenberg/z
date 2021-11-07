package writer

import (
	"encoding/json"
	"io"

	"github.com/friedenberg/z/lib/zettel"
	"github.com/friedenberg/z/util/stdprinter"
)

type Json struct {
	IncludeBody bool
}

func (p *Json) WriteZettel(w io.Writer, i int, z *zettel.Zettel) {
	var out interface{}

	if p.IncludeBody {
		out = map[string]interface{}{
			"metadata": z.Metadata,
			"body":     z.Body,
		}
	} else {
		out = z.Metadata
	}

	b, errOut := json.Marshal(out)

	if errOut != nil {
		stdprinter.Err(errOut)
		return
	}

	_, err := w.Write(b)
	stdprinter.PanicIfError(err)
}
