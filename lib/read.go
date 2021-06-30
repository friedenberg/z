package lib

import "fmt"

func (z *Zettel) HydrateFromFilePath() (err error) {
	err = z.ReadMetadata()

	if err != nil {
		err = fmt.Errorf("reading metadata: %w", err)
		return
	}

	err = z.ParseMetadata()

	if err != nil {
		err = fmt.Errorf("reading parsing: %w", err)
		return
	}

	return
}

func (z *Zettel) GenerateAlfredItemData() (err error) {
	err = z.AddAlfredItem()

	if err != nil {
		err = fmt.Errorf("adding alfred item: %w", err)
		return
	}

	err = z.GenerateAlfredJson()

	if err != nil {
		err = fmt.Errorf("generating alfred json: %w", err)
		return
	}

	return
}
