package util

import (
	"fmt"
	"os/exec"
)

type Git struct {
	Path string
}

type GitFilesToCommit struct {
	Git
	Files []string
}

func (g Git) CheckDiff() (ok bool, err error) {
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

func (g GitFilesToCommit) Add() (err error) {
	cmd := exec.Command(
		"git",
		append([]string{"add"}, g.Files...)...,
	)

	cmd.Dir = g.Path

	o, err := cmd.CombinedOutput()

	if err != nil {
		err = fmt.Errorf("git add: %w: %s", err, o)
		return
	}

	return
}

func (g GitFilesToCommit) Commit(msg string) (err error) {
	cmd := exec.Command(
		"git",
		"commit",
		"-S",
		"-m",
		msg,
	)

	cmd.Dir = g.Path
	o, err := cmd.CombinedOutput()

	if err != nil {
		err = fmt.Errorf("git commit: %w: %s", err, o)
		return
	}

	return
}

func (g GitFilesToCommit) AddAndCommit(msg string) (err error) {
	err = g.Add()

	if err != nil {
		return
	}

	ok := false

	if ok, err = g.CheckDiff(); !ok {
		return
	}

	err = g.Commit(msg)

	if err != nil {
		return
	}

	return
}
