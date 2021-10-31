package lib

import "github.com/friedenberg/z/lib/zettel"

func (k *FileStore) transactionProcessAdd(u Umwelt, z *Zettel) (err error) {
	if z.Id == 0 {
		var id zettel.Id
		id, err = u.Kasten.NewId()

		if err != nil {
			return
		}

		z.Id = id.Int()
		z.Path = MakePathFromId(u.Kasten.BasePath(), id.String())
	}

	err = k.readAndWrite(z, true)

	if err != nil {
		return
	}

	return
}
