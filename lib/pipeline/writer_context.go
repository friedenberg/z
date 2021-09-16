package pipeline

import (
	"io"
	"sync"

	"github.com/friedenberg/z/lib"
)

type WriterContext struct {
	i int
	z *lib.Zettel
	sync.Locker
	io.Writer
}

func makeWriterContext(i int, z *lib.Zettel, l sync.Locker, w io.Writer) (wc WriterContext) {
	wc.i = i
	wc.z = z
	wc.Locker = l
	wc.Writer = w
	return
}
