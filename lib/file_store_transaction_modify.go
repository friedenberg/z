package lib

import "github.com/friedenberg/z/lib/zettel"

func (k *FileStore) transactionProcessModify(u Umwelt, z *zettel.Zettel) (err error) {
	err = k.hydrateFromFileIfExists(z)

	if err != nil {
		return
	}

	err = k.writeToFile(z)

	if err != nil {
		return
	}

	return
}
