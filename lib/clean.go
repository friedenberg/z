package lib

import (
	"os"

	"github.com/friedenberg/z/util"
)

type CleanActionCheck func(z *Zettel) bool
type CleanActionPerform func(z *Zettel) error

type CleanAction struct {
	Check   CleanActionCheck
	Perform CleanActionPerform
}

func GetCleanActions() map[string]CleanAction {
	return map[string]CleanAction{
		"delete_if_missing_file": CleanAction{shouldDeleteIfMissingFile, deleteIfMissingFile},
		"add_date":               CleanAction{isMissingDate, addDate},
		// "rewrite_metadata": CleanAction{
		// 	func(_ *Zettel) bool { return true },
		// 	func(z *Zettel) error {
		// 		util.OpenFilesGuardInstance.Lock()
		// 		defer util.OpenFilesGuardInstance.Unlock()
		// 		return z.Write(nil)
		// 	},
		// },
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

	return !util.FileExists(z.IndexData.File)
}

func deleteIfMissingFile(z *Zettel) error {
	return os.Remove(z.Path)
}

func isMissingDate(z *Zettel) bool {
	return z.IndexData.Date == ""
}

func addDate(z *Zettel) (err error) {
	t, err := TimeFromPath(z.Path)

	if err != nil {
		return
	}

	z.IndexData.Date = t.Format("2006-01-02")

	err = z.Write(func(_ *Zettel, _ error) error { return nil })

	return
}
