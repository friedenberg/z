package lib

import (
	"github.com/friedenberg/z/util"
)

type Transactor func(Umwelt, Transaction) error

func (u Umwelt) RunTransaction(f Transactor) (err error) {
	t := Transaction{
		Add: &transactionPrinter{},
		Mod: &transactionPrinter{},
		Del: &transactionPrinter{},
	}

	f(u, t)

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

	for _, z := range t.Added() {
		z.Hydrate(true)
		u.Index.Add(z)
	}

	for _, z := range t.Modified() {
		z.Hydrate(true)
		u.Index.Update(z)
	}

	for _, z := range t.Deleted() {
		z.ReadMetadataAndBody()
		u.Index.Delete(z)
	}

	err = u.CacheIndex()

	if err != nil {
		return
	}

	return
}
