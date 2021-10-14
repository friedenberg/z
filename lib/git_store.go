package lib

import (
	"github.com/friedenberg/z/util/git"
)

type GitStore struct {
	FileStore
	SignCommits bool
}

func (k GitStore) CommitTransaction(u Umwelt) (err error) {
	err = k.FileStore.CommitTransaction(u)

	if err != nil {
		return
	}

	if u.Transaction.ShouldSkipCommit {
		return
	}

	//TODO-P2 migrate to git_store
	g := git.FilesToCommit{
		Git: git.Git{
			Path:       u.Kasten.Local.BasePath(),
			SignOption: git.SignOption(k.SignCommits),
		},
	}

	fileListMap := map[string][]string{
		"delete": u.Transaction.Deleted().Paths(),
		"modify": u.Transaction.Modified().Paths(),
		"add":    u.Transaction.Added().Paths(),
	}

	for k, v := range fileListMap {
		g.AddedOrModifiedFiles = nil
		g.DeletedFiles = nil

		switch k {
		case "modify":
			fallthrough
		case "add":
			g.AddedOrModifiedFiles = v
			err = g.Add()

		case "delete":
			g.DeletedFiles = v
			err = g.Delete()
		}

		if err != nil {
			return
		}

		ok := false

		if ok, err = g.HasChangesInDiff(); !ok {
			continue
		}

		if err != nil {
			return
		}

		err = g.Commit(k)

		if err != nil {
			return
		}
	}

	return
}
