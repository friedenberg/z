package lib

import (
	"os"

	"github.com/friedenberg/z/util"
)

type Transactor func(Umwelt) error

func (u Umwelt) RunTransaction(f Transactor) (err error) {
	f(u)

	readAndWrite := func(z *Zettel) (err error) {
		err = z.Hydrate(true)

		if os.IsNotExist(err) {
			err = nil
		} else if err != nil {
			return
		}

		err = z.Write(nil)

		if err != nil {
			return
		}

		return
	}

	for _, z := range u.Transaction.Added() {
		err = readAndWrite(z)

		if err != nil {
			return
		}

		u.Index.Add(z)
	}

	for _, z := range u.Transaction.Modified() {
		err = readAndWrite(z)

		if err != nil {
			return
		}

		u.Index.Update(z)
	}

	err = u.gitCommitTransactionIfNecessary()

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

func (u Umwelt) gitCommitTransactionIfNecessary() (err error) {
	if u.Transaction.ShouldSkipCommit {
		return
	}

	git := util.GitFilesToCommit{
		Git: util.Git{
			Path: u.Kasten.Local.BasePath,
		},
	}

	fileListMap := map[string][]string{
		"delete": u.Transaction.Deleted().Paths(),
		"modify": u.Transaction.Modified().Paths(),
		"add":    u.Transaction.Added().Paths(),
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
