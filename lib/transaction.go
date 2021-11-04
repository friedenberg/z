package lib

func MakeTransaction() (t Transaction) {
	t = Transaction{
		Add: MakeSet(),
		Mod: MakeSet(),
		Del: MakeSet(),
	}

	return
}

//TODO-P1 handle cases where files or zettles are just opened but not edited
type Transaction struct {
	ShouldSkipCommit bool
	ShouldCopyFiles  bool
	Add              Set
	Mod              Set
	Del              Set
}

func (t Transaction) RawFiles() (f []string) {
	f = make([]string, 0, t.Add.Len()+t.Mod.Len()+t.Del.Len())

	add := func(s ZettelSlice) {
		for _, z := range s {
			f = append(f, z.Path)
		}
	}

	add(t.Add.Zettels())
	add(t.Mod.Zettels())
	add(t.Del.Zettels())

	return
}
