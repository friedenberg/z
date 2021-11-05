package modifier

import "github.com/friedenberg/z/lib/zettel"

type modifier struct {
	modifierFunc ModifierFunc
}

func (m modifier) ModifyZettel(i int, z *zettel.Zettel) (err error) {
	if m.modifierFunc != nil {
		err = m.modifierFunc(i, z)
	}

	return
}
