package lib

type Transactor func(Umwelt, Transaction) error

func (u Umwelt) RunTransaction(f Transactor) (err error) {
	t := Transaction{}
	f(u, t)
	//TODO process
	return
}
