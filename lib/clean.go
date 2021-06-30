package lib

import "os"

type CleanActionCheck func(z *Zettel) bool
type CleanActionPerform func(z *Zettel) error

type CleanAction struct {
	Check   CleanActionCheck
	Perform CleanActionPerform
}

func GetCleanActions() map[string]CleanAction {
	return map[string]CleanAction{
		"delete_if_missing_file": CleanAction{
			shouldDeleteIfMissingFile,
			deleteIfMissingFile,
		},
		//TODO change file permissions
		//TODO reformat yaml
	}
}

func shouldDeleteIfMissingFile(z *Zettel) bool {
	if z.Metadata.Kind != "file" {
		return false
	}

	_, err := os.Stat(z.Metadata.File)
	return os.IsNotExist(err)
}

func deleteIfMissingFile(z *Zettel) error {
	return os.Remove(z.Path)
}
