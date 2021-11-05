package modifier

import "github.com/friedenberg/z/lib"

func TransactionAction(t lib.Transaction, a lib.TransactionAction) modifier {
	return Make(
		func(i int, z *lib.Zettel) (err error) {
			t.Set(z, a)
			return
		},
	)
}
