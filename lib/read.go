package lib

import "path/filepath"

func (z *Zettel) HydrateFromFilePath(path string) (err error) {
	z.Path, err = filepath.Abs(path)

	if err != nil {
		return
	}

	err = z.ReadMetadata()

	if err != nil {
		return
	}

	err = z.ParseMetadata()

	return
}

func (z *Zettel) GenerateAlfredItemData() (err error) {
	err = z.AddAlfredItem()

	if err != nil {
		return
	}

	err = z.GenerateAlfredJson()

	return
}
