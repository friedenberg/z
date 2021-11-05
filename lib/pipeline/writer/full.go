package writer

import (
	"io"
	"strings"

	"github.com/friedenberg/z/lib/zettel"
	"github.com/friedenberg/z/util/stdprinter"
)

type Full struct{}

func (p *Full) WriteZettel(w io.Writer, _ int, z *zettel.Zettel) {
	sb := &strings.Builder{}

	y, err := z.Note.Metadata.ToYAMLWithBoundary()

	if err != nil {
		stdprinter.Error(err)
		return
	}

	sb.WriteString(y)
	sb.WriteString(z.Body)
	_, err = io.WriteString(w, sb.String())
	stdprinter.PanicIfError(err)
}
