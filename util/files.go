package util

import (
	"os"
	"path"
)

func FileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func BaseNameNoSuffix(p string) string {
	b := path.Base(p)
	return b[0 : len(b)-len(path.Ext(b))]
}
