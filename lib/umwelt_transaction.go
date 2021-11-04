package lib

type Transactor func(Umwelt) error

func (u Umwelt) RunTransaction(f Transactor) (err error) {
	if f != nil {
		err = f(u)
	}

	if err != nil {
		return
	}

	err = u.Kasten.CommitTransaction(u)

	if err != nil {
		return
	}

	err = u.IndexTransaction()

	if err != nil {
		return
	}

	return
}

func (u Umwelt) IndexTransaction() (err error) {
	for _, z := range u.Add.Zettels() {
		u.Index.Add(z)
	}

	for _, z := range u.Mod.Zettels() {
		u.Index.Update(z)
	}

	for _, z := range u.Del.Zettels() {
		u.Index.Delete(z)
	}

	err = u.CacheIndex()

	if err != nil {
		return
	}

	return
}
