package git

import (
	"os/exec"

	"github.com/friedenberg/z/util"
	"golang.org/x/xerrors"
)

type SignOption bool

func (so SignOption) String() string {
	if so {
		return "--gpg-sign"
	} else {
		return "--no-gpg-sign"
	}
}

type Git struct {
	Path string
	SignOption
}

type FilesToCommit struct {
	Git
	AddedOrModifiedFiles []string
	DeletedFiles         []string
}

func (g Git) HasChangesInDiff() (ok bool, err error) {
	cmd := exec.Command(
		"git",
		"diff",
		"--staged",
		"--quiet",
	)

	cmd.Dir = g.Path

	err = cmd.Run()

	if _, ok = err.(*exec.ExitError); ok {
		err = nil
	} else if err == nil {
		ok = false
		err = nil
	}

	return
}

func (g FilesToCommit) Add() (err error) {
	if len(g.AddedOrModifiedFiles) == 0 {
		return
	}

	cmd := util.ExecCommand(
		"git",
		[]string{"add"},
		g.AddedOrModifiedFiles,
	)

	cmd.Dir = g.Path

	o, err := cmd.CombinedOutput()

	if err != nil {
		err = xerrors.Errorf("git add: %w: %s", err, o)
		return
	}

	return
}

func (g FilesToCommit) Delete() (err error) {
	if len(g.DeletedFiles) == 0 {
		return
	}

	cmd := util.ExecCommand(
		"git",
		[]string{"rm", "--ignore-unmatch"},
		g.DeletedFiles,
	)

	cmd.Dir = g.Path

	o, err := cmd.CombinedOutput()

	if err != nil {
		err = xerrors.Errorf("git add: %w: %s", err, o)
		return
	}

	return
}

func (g FilesToCommit) Commit(msg string) (err error) {
	cmd := exec.Command(
		"git",
		"commit",
		g.SignOption.String(),
		"-m",
		msg,
	)

	cmd.Dir = g.Path
	o, err := cmd.CombinedOutput()

	if err != nil {
		err = xerrors.Errorf("git commit: %w: %s", err, o)
		return
	}

	return
}

func (g FilesToCommit) AddAndDeleteAndCommit(msg string) (err error) {
	err = g.Delete()

	if err != nil {
		return
	}

	err = g.Add()

	if err != nil {
		return
	}

	ok := false

	if ok, err = g.HasChangesInDiff(); !ok {
		return
	}

	err = g.Commit(msg)

	if err != nil {
		return
	}

	return
}
