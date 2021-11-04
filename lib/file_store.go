package lib

import (
	"bufio"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/friedenberg/z/lib/zettel"
	"github.com/friedenberg/z/util/files_guard"
	"github.com/friedenberg/z/util/stdprinter"
	"golang.org/x/xerrors"
)

var durationOneSecond time.Duration

func init() {
	var err error
	durationOneSecond, err = time.ParseDuration("1s")
	stdprinter.PanicIfError(err)
}

//TODO-P2 move to lib/kasten
type FileStore struct {
	umwelt     Umwelt
	basePath   string
	lastIdTime time.Time
	*sync.Mutex
}

func (s *FileStore) Init(u Umwelt, o map[string]interface{}) (err error) {
	s.umwelt = u
	//TODO-P1 init from options
	s.lastIdTime = time.Now()
	s.Mutex = &sync.Mutex{}
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
	id = k.nextAvailableId()
	return
}

func (k *FileStore) nextAvailableId() (id zettel.Id) {
	k.Lock()
	defer k.Unlock()

	k.lastIdTime = k.lastIdTime.Add(durationOneSecond)
	id = zettel.Id(k.lastIdTime.Unix())

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
	for _, z := range u.Transaction.Add.Zettels() {
		err = k.transactionProcessAdd(u, z)

		if err != nil {
			return
		}
	}

	for _, z := range u.Transaction.Mod.Zettels() {
		err = k.transactionProcessModify(u, z)

		if err != nil {
			return
		}
	}

	for _, z := range u.Transaction.Del.Zettels() {
		err = k.transactionProcessDelete(u, z)

		if err != nil {
			return
		}
	}

	return
}
