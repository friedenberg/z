package pipeline

import (
	"io"
	"sync"

	"github.com/friedenberg/z/lib/zettel"
)

type WriterContext struct {
	i int
	z *zettel.Zettel
	sync.Locker
	io.Writer
}

func makeWriterContext(i int, z *zettel.Zettel, l sync.Locker, w io.Writer) (wc WriterContext) {
	wc.i = i
	wc.z = z
	wc.Locker = l
	wc.Writer = w
	return
}
