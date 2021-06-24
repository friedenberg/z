package lib

import "os"

type CleanActionCheck func() bool
type CleanActionPerform func() error

type CleanAction struct {
	Check   CleanActionCheck
	Perform CleanActionPerform
}

func (z *Zettel) GetCleanActions() map[string]CleanAction {
	return map[string]CleanAction{
		"delete_if_missing_file": CleanAction{
			z.ShouldDeleteIfMissingFile,
			z.DeleteIfMissingFile,
		},
	}
}

func (z *Zettel) ShouldDeleteIfMissingFile() bool {
	if z.Metadata.Kind != "file" {
		return false
	}

	_, err := os.Stat(z.Metadata.File)
	return os.IsNotExist(err)
}

func (z *Zettel) DeleteIfMissingFile() error {
	return os.Remove(z.Path)
}
