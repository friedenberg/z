package lib

import (
	"github.com/friedenberg/z/util"
)

type Transactor func(Umwelt, *Transaction) error

func (u Umwelt) RunTransaction(f Transactor) (err error) {
	t := &Transaction{
		Add: &transactionPrinter{},
		Mod: &transactionPrinter{},
		Del: &transactionPrinter{},
	}

	f(u, t)

	for _, z := range t.Added() {
		err = z.Write(nil)

		if err != nil {
			return
		}

		u.Index.Add(z)
	}

	for _, z := range t.Modified() {
		err = z.Write(nil)

		if err != nil {
			return
		}

		u.Index.Update(z)
	}

	err = u.gitCommitTransactionIfNecessary(t)

	if err != nil {
		return
	}

	for _, z := range t.Added() {
		u.Index.Add(z)
	}

	for _, z := range t.Modified() {
		u.Index.Update(z)
	}

	for _, z := range t.Deleted() {
		u.Index.Delete(z)
	}

	err = u.CacheIndex()

	if err != nil {
		return
	}

	return
}

func (u Umwelt) gitCommitTransactionIfNecessary(t *Transaction) (err error) {
	if t.ShouldSkipCommit {
		return
	}

	git := util.GitFilesToCommit{
		Git: util.Git{
			Path: u.Kasten.Local.BasePath,
		},
	}

	fileListMap := map[string][]string{
		"delete": t.Deleted().Paths(),
		"modify": t.Modified().Paths(),
		"add":    t.Added().Paths(),
	}

	for k, v := range fileListMap {
		git.Files = v
		err = git.AddAndCommit(k)

		if err != nil {
			return
		}
	}

	return
}
