package writer

import (
	"encoding/json"
	"io"

	"github.com/friedenberg/z/lib"
	"github.com/friedenberg/z/util/stdprinter"
	"github.com/pelletier/go-toml/v2"
	"golang.org/x/xerrors"
)

type TomlToJson struct{}

func (p *TomlToJson) WriteZettel(w io.Writer, i int, z *lib.Zettel) {
	if _, ok := z.Metadata.TagSet().Get("k-toml"); !ok {
		return
	}

	doc, errOut := TomlZettelBody(z)

	if errOut != nil {
		errOut = xerrors.Errorf("failed to unmarshal for %d: %w", z.Id, errOut)
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

	_, err := w.Write(b)
	stdprinter.PanicIfError(err)
}

func TomlZettelBody(z *lib.Zettel) (out map[string]interface{}, err error) {
	err = toml.Unmarshal([]byte(z.Body), &out)
	return
}
