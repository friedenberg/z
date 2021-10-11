package lib

import (
	"bufio"
	"bytes"
	"path"
	"strconv"
	"strings"

	"github.com/friedenberg/z/lib/zettel/metadata"
	"github.com/friedenberg/z/util/files_guard"
	"golang.org/x/xerrors"
)

func (z *Zettel) Hydrate(readBody bool) (err error) {
	id := strings.TrimSuffix(path.Base(z.Path), path.Ext(z.Path))
	idInt, err := strconv.ParseInt(id, 10, 64)

	if err != nil {
		err = xerrors.Errorf("extracting id from filename: %w", err)
		return
	}

	z.Id = idInt

	f, err := files_guard.Open(z.Path)
	defer files_guard.Close(f)

	if err != nil {
		return
	}

	r := bufio.NewReader(f)

	y, err := metadata.ReadYAMLHeader(r)

	if err != nil {
		return
	}

	err = z.Metadata.Set(y)

	if err != nil {
		return
	}

	if readBody {
		body := &bytes.Buffer{}
		_, err = r.WriteTo(body)

		if err != nil {
			return
		}

		z.Body = string(body.Bytes())
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
