package lib

import (
	"os"
	"path"

	"github.com/friedenberg/z/util"
)

//TODO swithch to p rintable description
type CleanActionCheck func(z *Zettel) bool
type CleanActionPerform func(z *Zettel) (bool, error)

type CleanAction struct {
	Check   CleanActionCheck
	Perform CleanActionPerform
}

func GetCleanActions() map[string]CleanAction {
	return map[string]CleanAction{
		"delete_if_missing_file": CleanAction{shouldDeleteIfMissingFile, deleteIfMissingFile},
		"normalize_file": CleanAction{
			func(z *Zettel) bool {
				if z.IndexData.File == "" {
					return false
				}

				normalizedFile := path.Base(z.IndexData.File)

				return normalizedFile != z.IndexData.File
			},
			func(z *Zettel) (shouldWrite bool, err error) {
				z.IndexData.File = path.Base(z.IndexData.File)
				shouldWrite = true
				return
			},
		},
		"rewrite_metadata": CleanAction{
			func(z *Zettel) bool {
				oldYaml := z.Data.MetadataYaml
				//TODO handle err
				z.GenerateMetadataYaml()

				return oldYaml != z.Data.MetadataYaml
			},
			func(z *Zettel) (shouldWrite bool, err error) {
				util.OpenFilesGuardInstance.Lock()
				defer util.OpenFilesGuardInstance.Unlock()
				shouldWrite = true
				return
			},
		},
		//index
		// "remove_from_index":             CleanAction{shouldRemoveFromIndex, removeFromIndex},
		// "add_to_index":             CleanAction{shouldRemoveFromIndex, removeFromIndex},
		// "update_in_index":             CleanAction{shouldUpdateInIndex, updateInIndex},
		//TODO file attachment
		//TODO change file permissions
		//TODO reformat yaml
		//TODO add date
	}
}

func shouldDeleteIfMissingFile(z *Zettel) bool {
	if !z.HasFile() {
		return false
	}

	return !util.FileExists(z.FilePath())
}

func deleteIfMissingFile(z *Zettel) (shouldWrite bool, err error) {
	err = os.Remove(z.Path)
	return
}
