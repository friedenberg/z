package modifier

import (
	"github.com/friedenberg/z/lib"
	"github.com/friedenberg/z/lib/zettel"
)

func TransactionAction(t lib.Transaction, a lib.TransactionAction) modifier {
	return Make(
		func(i int, z *zettel.Zettel) (err error) {
			t.Set(z, a)
			return
		},
	)
}
