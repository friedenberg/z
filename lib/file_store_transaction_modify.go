package lib

func (k *FileStore) transactionProcessModify(u Umwelt, z *Zettel) (err error) {
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
