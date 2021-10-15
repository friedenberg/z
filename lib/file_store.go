package lib

import (
	"bufio"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/friedenberg/z/lib/zettel"
	"github.com/friedenberg/z/lib/zettel/metadata"
	"github.com/friedenberg/z/util"
	"github.com/friedenberg/z/util/files_guard"
	"golang.org/x/xerrors"
)

//TODO-P2 move to lib/kasten
type FileStore struct {
	umwelt   Umwelt
	basePath string
}

func (s *FileStore) Init(u Umwelt, o map[string]interface{}) (err error) {
	s.umwelt = u
	//TODO-P1 init from options
	return
}

func (s FileStore) BasePath() string {
	return s.basePath
}

func (e FileStore) GetAll() (zettels []string, err error) {
	glob := filepath.Join(e.BasePath(), "*.md")
	zettels, err = filepath.Glob(glob)
	return
}

func (e FileStore) GetNormalizedPath(a string) (b string, err error) {
	if filepath.IsAbs(a) {
		b = a
	} else {
		b, err = filepath.Abs(path.Join(e.BasePath(), a))
	}

	return
}

func (k *FileStore) NewId() (id zettel.Id, err error) {
	currentTime := time.Now()

	for {
		p := MakePathFromTime(k.BasePath(), currentTime)

		if util.FileExists(p) {
			d, err := time.ParseDuration("1s")

			if err != nil {
				panic(err)
			}

			currentTime = currentTime.Add(d)
		} else {
			id = zettel.Id(currentTime.Unix())
			return
		}
	}

	return
}

func (k FileStore) Hydrate(z *Zettel, includeBody bool) (err error) {
	z.Umwelt = k.umwelt

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

	z.ReadFrom(r, includeBody)

	if err != nil {
		return
	}

	return
}

func (k FileStore) CommitTransaction(u Umwelt) (err error) {
	readAndWrite := func(z *Zettel) (err error) {
		err = k.Hydrate(z, true)

		if os.IsNotExist(err) {
			err = nil
		} else if err != nil {
			return
		}

		err = z.Write(nil)

		if err != nil {
			return
		}

		return
	}

	for _, z := range u.Transaction.Added() {
		err = readAndWrite(z)

		if err != nil {
			return
		}
	}

	for _, z := range u.Transaction.Modified() {
		err = readAndWrite(z)

		if err != nil {
			return
		}
	}

	return
}

func (k FileStore) CopyFileTo(localPath string, fd metadata.File) (err error) {
	remotePath := path.Join(k.BasePath(), fd.FileName())

	cmd := exec.Command("cp", "-R", localPath, remotePath)
	out, err := cmd.CombinedOutput()

	if err != nil {
		err = xerrors.Errorf("%w: %s", err, out)
		return
	}

	return
}

func (k FileStore) CopyFileFrom(localPath string, fd metadata.File) (err error) {
	remotePath := path.Join(k.BasePath(), fd.FileName())

	cmd := exec.Command("cp", "-R", remotePath, localPath)
	out, err := cmd.CombinedOutput()

	if err != nil {
		err = xerrors.Errorf("%w: %s", err, out)
		return
	}

	return
}
