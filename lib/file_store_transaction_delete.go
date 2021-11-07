package lib

import (
	"os"

	"github.com/friedenberg/z/lib/zettel"
	"github.com/friedenberg/z/util"
)

func (k *FileStore) transactionProcessDelete(u *Umwelt, z *zettel.Zettel) (err error) {
	err = k.hydrateFromFileIfExists(z)

	if err != nil {
		return
	}

	err = os.Remove(z.Path)

	if err != nil {
		return
	}

	if f := z.Metadata.File(); f != nil {
		path := f.FilePath(k.umwelt.Dir())
		err = util.SetAllowUserChanges(path)

		if err != nil {
			return
		}

		err = os.Remove(f.FilePath(k.umwelt.Dir()))

		if err != nil {
			return
		}
	}

	u.Set(z, TransactionActionDeleted)

	return
}
