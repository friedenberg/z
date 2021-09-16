package modifier

import "github.com/friedenberg/z/lib"

type modifier struct {
	modifierFunc ModifierFunc
}

func (m modifier) ModifyZettel(i int, z *lib.Zettel) (err error) {
	if m.modifierFunc != nil {
		err = m.modifierFunc(i, z)
	}

	return
}
