package lib

import (
	"os"
)

func (k *FileStore) hydrateFromFileIfExists(z *Zettel) (err error) {
	err = k.Hydrate(z, true)

	if os.IsNotExist(err) {
		err = nil
	} else if err != nil {
		return
	}

	return
}

func (k *FileStore) writeToFile(z *Zettel) (err error) {
	err = z.Write(nil)

	if err != nil {
		return
	}

	return
}
