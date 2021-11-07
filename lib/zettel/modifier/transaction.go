package modifier

import (
	"github.com/friedenberg/z/lib"
	"github.com/friedenberg/z/lib/zettel"
	"github.com/friedenberg/z/util/stdprinter"
)

func TransactionAction(t lib.Transaction, a lib.TransactionAction) modifier {
	return Make(
		func(i int, z *zettel.Zettel) (err error) {
			stdprinter.Debug(
				"TransactionAction",
				"setting zettel",
				a,
				z.Path,
			)
			t.Set(z, a)
			return
		},
	)
}
