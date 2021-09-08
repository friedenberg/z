package util

import (
	"fmt"
)

type GitAnnex struct {
	GitFilesToCommit
}

func (g GitAnnex) Unlock() (err error) {
	cmd := ExecCommand(
		"git",
		[]string{"annex", "unlock"},
		g.Files,
	)

	cmd.Dir = g.Path

	o, err := cmd.CombinedOutput()

	if err != nil {
		err = fmt.Errorf("git annex unlock: %w: %s", err, o)
		return
	}

	return
}

func (g GitAnnex) Lock() (err error) {
	cmd := ExecCommand(
		"git",
		[]string{"annex", "lock"},
		g.Files,
	)

	cmd.Dir = g.Path

	o, err := cmd.CombinedOutput()

	if err != nil {
		err = fmt.Errorf("git annex unlock: %w: %s", err, o)
		return
	}

	return
}
