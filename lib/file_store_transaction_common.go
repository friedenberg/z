package lib

import (
	"os"

	"github.com/friedenberg/z/lib/zettel"
)

func (k *FileStore) hydrateFromFileIfExists(z *zettel.Zettel) (err error) {
	err = k.Hydrate(z, true)

	if os.IsNotExist(err) {
		err = nil
	} else if err != nil {
		return
	}

	return
}

func (k *FileStore) writeToFile(z *zettel.Zettel) (err error) {
	err = z.Write(nil)

	if err != nil {
		return
	}

	return
}
