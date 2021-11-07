package lib

import (
	"github.com/friedenberg/z/lib/zettel"
	"github.com/friedenberg/z/util/stdprinter"
)

func (k *FileStore) transactionProcessAdd(u *Umwelt, z *zettel.Zettel) (err error) {
	stdprinter.Debug("will process transaction add for zettel:", z.Path)
	//add can be called for existing zettels or new zettels
	//in the case of new, we need to create an id and populate it
	if z.Id == 0 {
		var id zettel.Id
		id, err = k.umwelt.Kasten.NewId()

		if err != nil {
			return
		}

		z.Id = id.Int()
		z.Path = MakePathFromId(k.umwelt.Kasten.BasePath(), id.String())
	}

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

	stdprinter.Debug("did process transaction add for zettel:", z.Path)
	return
}
