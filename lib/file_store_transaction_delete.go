package lib

import "os"

func (k *FileStore) transactionProcessDelete(u Umwelt, z *Zettel) (err error) {
	err = k.hydrateFromFileIfExists(z)

	if err != nil {
		return
	}

	err = os.Remove(z.Path)

	if err != nil {
		return
	}

	if f := z.Metadata.File(); f != nil {
		err = os.Remove(f.FilePath(u.BasePath))
	}

	if err != nil {
		return
	}

	u.Set(z, TransactionActionDeleted)

	return
}
