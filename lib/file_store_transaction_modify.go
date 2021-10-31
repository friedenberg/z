package lib

func (k *FileStore) transactionProcessModify(u Umwelt, z *Zettel) (err error) {
	err = k.readAndWrite(z, true)

	if err != nil {
		return
	}

	return
}
