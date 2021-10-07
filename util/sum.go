package util

import (
	"crypto/sha256"
	"fmt"
	"io"

	"github.com/friedenberg/z/util/files_guard"
)

func Sha256HashForFile(filepath string) (sum string, err error) {
	file, err := files_guard.Open(filepath)

	if err != nil {
		return
	}

	defer files_guard.Close(file)

	hash := sha256.New()

	if _, err = io.Copy(hash, file); err != nil {
		return
	}

	sum = fmt.Sprintf("%x", hash.Sum(nil))

	return
}
