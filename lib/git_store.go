package lib

import (
	"github.com/friedenberg/z/util/git"
)

type GitStore struct {
	FileStore
	SignCommits bool
}

func (k GitStore) CommitTransaction(u *Umwelt) (err error) {
	err = k.FileStore.CommitTransaction(u)

	if err != nil {
		return
	}

	if u.ShouldSkipCommit {
		return
	}

	g := git.FilesToCommit{
		Git: git.Git{
			Path:       k.BasePath(),
			SignOption: git.SignOption(k.SignCommits),
		},
	}

	run := func(k string, v []string) (err error) {
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
			return
		}

		if err != nil {
			return
		}

		err = g.Commit(k)

		if err != nil {
			return
		}

		return
	}

	err = run("add", u.ZettelsForActions(TransactionActionAdded).Paths())

	if err != nil {
		return
	}

	err = run("modify", u.ZettelsForActions(TransactionActionModified).Paths())

	if err != nil {
		return
	}

	err = run("delete", u.ZettelsForActions(TransactionActionDeleted).Paths())

	if err != nil {
		return
	}

	return
}
