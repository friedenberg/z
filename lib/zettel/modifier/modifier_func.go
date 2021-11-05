package modifier

import "github.com/friedenberg/z/lib/zettel"

type ModifierFunc func(int, *zettel.Zettel) error
