package util

import (
	"os"
	"path"
	"testing"
)

const expected = "8ee5fe0522cecebe0861eda68f19c118f3582e7ad14321069f39f5d9110930c6"

func TestSum(t *testing.T) {
	cwd, err := os.Getwd()

	if err != nil {
		t.Fatalf(err.Error())
	}

	path := path.Join(cwd, "sum.go")
	hash, err := Sha256HashForFile(path)

	if err != nil {
		t.Fatal(err)
	}

	if hash != expected {
		t.Errorf("expected '%s' for hash but got '%s'", expected, hash)
	}
}
