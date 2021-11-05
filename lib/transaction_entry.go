package lib

import "github.com/friedenberg/z/lib/zettel"

type TransactionEntry struct {
	*zettel.Zettel
	TransactionAction
}
