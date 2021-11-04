package util

import (
	"os"
	"os/exec"
	"path"

	"golang.org/x/xerrors"
)

func setUserChanges(path string, allow bool) (err error) {
	setting := "uchg"

	if allow {
		setting = "no" + setting
	}

	cmd := exec.Command("chflags", setting, path)
	var msg []byte
	msg, err = cmd.CombinedOutput()

	if err != nil {
		err = xerrors.Errorf("failed to run chflags: %s", msg)
		return
	}
	return
}

func SetAllowUserChanges(path string) (err error) {
	return setUserChanges(path, true)
}

func SetDisallowUserChanges(path string) (err error) {
	return setUserChanges(path, false)
}

func FileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func IsDir(path string) (isDir bool, err error) {
	s, err := os.Stat(path)

	if os.IsNotExist(err) {
		err = nil
	} else if err == nil && s.Mode().IsDir() {
		isDir = true
	}

	return
}

func BaseNameNoSuffix(p string) string {
	b := path.Base(p)
	return b[0 : len(b)-len(path.Ext(b))]
}

func ExtNoDot(p string) (e string) {
	e = path.Ext(p)

	if len(e) != 0 {
		e = e[1:]
	}

	return
}

func EverythingExceptExtension(p string) string {
	return p[0 : len(p)-len(path.Ext(p))]
}
