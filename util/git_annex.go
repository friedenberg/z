package util

import (
	"fmt"
	"os/exec"
)

type GitAnnex struct {
	GitFilesToCommit
}

func (g *GitAnnex) Unlock() (err error) {
	cmd := exec.Command(
		"git",
		append([]string{"annex", "unlock"}, g.Files...)...,
	)

	cmd.Dir = g.Path

	o, err := cmd.CombinedOutput()

	if err != nil {
		err = fmt.Errorf("git annex unlock: %w: %s", err, o)
		return
	}

	return
}

func (g *GitAnnex) Lock() (err error) {
	cmd := exec.Command(
		"git",
		append([]string{"annex", "lock"}, g.Files...)...,
	)

	cmd.Dir = g.Path

	o, err := cmd.CombinedOutput()

	if err != nil {
		err = fmt.Errorf("git annex unlock: %w: %s", err, o)
		return
	}

	return
}
