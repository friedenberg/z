package lib

import (
	"github.com/friedenberg/z/lib/zettel"
	"github.com/friedenberg/z/util/stdprinter"
)

func (k *FileStore) transactionProcessModify(u *Umwelt, z *zettel.Zettel) (err error) {
	stdprinter.Debug("will process transaction modify for zettel:", z.Path)
	err = k.hydrateFromFileIfExists(z)

	if err != nil {
		return
	}

	if u.IsFinalTransaction {
		err = k.updateFilesIfNecessary(z)

		if err != nil {
			return
		}
	}

	err = k.writeToFile(z)

	if err != nil {
		return
	}

	stdprinter.Debug("did process transaction modify for zettel:", z.Path)
	return
}
