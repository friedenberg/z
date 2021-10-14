package lib

type Transactor func(Umwelt) error

func (u Umwelt) RunTransaction(f Transactor) (err error) {
	f(u)

	u.Kasten.Local.CommitTransaction(u)

	if err != nil {
		return
	}

	for _, z := range u.Added() {
		u.Index.Add(z)
	}

	for _, z := range u.Modified() {
		u.Index.Update(z)
	}

	for _, z := range u.Deleted() {
		u.Index.Delete(z)
	}

	err = u.CacheIndex()

	if err != nil {
		return
	}

	return
}
