package lib

import (
	"github.com/friedenberg/z/lib/zettel"
	"github.com/friedenberg/z/lib/zettel/modifier"
	"github.com/friedenberg/z/util/stdprinter"
)

func MakeTransactionAction(t *Transaction, a TransactionAction) modifier.Modifier {
	return modifier.Make(
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
