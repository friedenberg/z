package lib

import "github.com/friedenberg/z/lib/kasten"

type Umwelt struct {
	DefaultKasten kasten.Implementation
	Kasten        map[string]kasten.Implementation
}

func (u Umwelt) FilesAndGit() *FilesAndGit {
	return u.DefaultKasten.(*FilesAndGit)
}
