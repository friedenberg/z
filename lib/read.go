package lib

import (
	"path"
	"strconv"
	"strings"

	"golang.org/x/xerrors"
)

func (z *KastenZettel) Hydrate(readBody bool) (err error) {
	id := strings.TrimSuffix(path.Base(z.Path), path.Ext(z.Path))
	idInt, err := strconv.ParseInt(id, 10, 64)

	if err != nil {
		err = xerrors.Errorf("extracting id from filename: %w", err)
		return
	}

	z.Id = idInt

	if readBody {
		err = z.ReadMetadataAndBody()
	} else {
		err = z.ReadMetadata()
	}

	if err != nil {
		err = xerrors.Errorf("reading metadata: %w", err)
		return
	}

	err = z.ParseMetadata()

	if err != nil {
		err = xerrors.Errorf("reading parsing: %w", err)
		return
	}

	return
}

// func (z *Zettel) GenerateAlfredItemData(f AlfredItemFormat) (err error) {
// 	err = z.AddAlfredItem(f)

// 	if err != nil {
// 		err = fmt.Errorf("adding alfred item: %w", err)
// 		return
// 	}

// 	err = z.GenerateAlfredJson()

// 	if err != nil {
// 		err = fmt.Errorf("generating alfred json: %w", err)
// 		return
// 	}

// 	return
// }
