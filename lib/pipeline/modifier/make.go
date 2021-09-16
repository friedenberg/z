package modifier

func Make(f ModifierFunc) (m modifier) {
	m.modifierFunc = f
	return
}
