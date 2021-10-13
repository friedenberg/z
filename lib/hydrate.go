package lib

import (
	"bufio"
	"path"
	"strconv"
	"strings"

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

	z.ReadFrom(r, readBody)

	if err != nil {
		return
	}

	return
}
