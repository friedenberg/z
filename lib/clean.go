package lib

import (
	"os"
	"path"

	"github.com/friedenberg/z/util"
)

//TODO swithch to p rintable description
type CleanActionCheck func(z *KastenZettel) bool
type CleanActionPerform func(z *KastenZettel) (bool, error)

type CleanAction struct {
	Check   CleanActionCheck
	Perform CleanActionPerform
}

func GetCleanActions() map[string]CleanAction {
	return map[string]CleanAction{
		"delete_if_missing_file": CleanAction{shouldDeleteIfMissingFile, deleteIfMissingFile},
		"normalize_file": CleanAction{
			func(z *KastenZettel) bool {
				if z.HasFile() {
					return false
				}

				normalizedFile := path.Base(z.Metadata.File)

				return normalizedFile != z.Metadata.File
			},
			func(z *KastenZettel) (shouldWrite bool, err error) {
				z.Metadata.File = path.Base(z.Metadata.File)
				shouldWrite = true
				return
			},
		},
		"rewrite_metadata": CleanAction{
			func(z *KastenZettel) bool {
				oldYaml := z.Data.MetadataYaml
				//TODO handle err
				z.generateMetadataYaml()

				return oldYaml != z.Data.MetadataYaml
			},
			func(z *KastenZettel) (shouldWrite bool, err error) {
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
	}
}

func shouldDeleteIfMissingFile(z *KastenZettel) bool {
	if !z.HasFile() {
		return false
	}

	return !util.FileExists(z.FilePath())
}

func deleteIfMissingFile(z *KastenZettel) (shouldWrite bool, err error) {
	err = os.Remove(z.Path)
	return
}
