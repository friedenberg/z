package lib

import (
	"os"
	"path/filepath"
	"strconv"
	"time"

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
		//TODO file attachment
		//TODO change file permissions
		//TODO reformat yaml
		//TODO add date
	}
}

func shouldDeleteIfMissingFile(z *Zettel) bool {
	if z.Metadata.Kind != "file" {
		return false
	}

	return !util.FileExists(z.Metadata.File)
}

func deleteIfMissingFile(z *Zettel) error {
	return os.Remove(z.Path)
}

func isMissingDate(z *Zettel) bool {
	return z.Metadata.Date == ""
}

func addDate(z *Zettel) (err error) {
	base := filepath.Base(z.Path)
	ext := filepath.Ext(z.Path)

	base = base[:len(base)-len(ext)]
	i, err := strconv.ParseInt(base, 10, 64)

	if err != nil {
		return
	}

	t := time.Unix(i, 0)
	z.Metadata.Date = t.Format("2006-01-02")

	err = z.Write()

	return
}
