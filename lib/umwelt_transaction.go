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
	for _, z := range u.ZettelsForActions(TransactionActionAdded) {
		u.Index.Add(z)
	}

	for _, z := range u.ZettelsForActions(TransactionActionModified) {
		u.Index.Update(z)
	}

	for _, z := range u.ZettelsForActions(TransactionActionDeleted) {
		u.Index.Delete(z)
	}

	err = u.CacheIndex()

	if err != nil {
		return
	}

	return
}
