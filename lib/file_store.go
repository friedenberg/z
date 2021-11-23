package lib

import (
	"bufio"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/friedenberg/z/lib/feeder"
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

type FileStore struct {
	umwelt     *Umwelt
	basePath   string
	lastIdTime time.Time
	*sync.Mutex
}

func (s *FileStore) Init(u *Umwelt, o map[string]interface{}) (err error) {
	s.umwelt = u
	s.lastIdTime = time.Now()
	s.Mutex = &sync.Mutex{}
	return
}

func (s FileStore) BasePath() string {
	return s.basePath
}

func (e FileStore) GetAll() feeder.Feeder {
	glob := filepath.Join(e.BasePath(), "*.md")
	zettels, err := filepath.Glob(glob)

	if err != nil {
		stdprinter.PanicIfError(err)
	}

	return feeder.MakeStringSlice(zettels)
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

func (k FileStore) Hydrate(u *Umwelt, z *zettel.Zettel, includeBody bool) (err error) {
	z.ZUmwelt = u

	id := strings.TrimSuffix(path.Base(z.Path), path.Ext(z.Path))
	idInt, err := strconv.ParseInt(id, 10, 64)

	if err != nil {
		err = xerrors.Errorf("extracting id from filename: %w", err)
		return
	}

	z.Id = zettel.Id(idInt)

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

func (k FileStore) CommitTransaction(u *Umwelt) (err error) {
	stdprinter.Debug("FileStore.CommitTransaction", "will commit transaction")

	if u.Len() == 0 {
		stdprinter.Debug("nothing to transact, terminating early")
		return
	}

	for _, z := range u.ZettelsForActions(TransactionActionAdded) {
		stdprinter.Debug(
			"FileStore.CommitTransaction",
			"will process add",
			z.Path,
		)

		err = k.transactionProcessAdd(u, z)

		if err != nil {
			return
		}
	}

	for _, z := range u.ZettelsForActions(TransactionActionModified) {
		stdprinter.Debug(
			"FileStore.CommitTransaction",
			"will process modify",
			z.Path,
		)

		err = k.transactionProcessModify(u, z)

		if err != nil {
			return
		}
	}

	for _, z := range u.ZettelsForActions(TransactionActionDeleted) {
		stdprinter.Debug(
			"FileStore.CommitTransaction",
			"will process delete",
			z.Path,
		)

		err = k.transactionProcessDelete(u, z)

		if err != nil {
			return
		}
	}

	for _, z := range u.ZettelsForActions(TransactionActionAdded) {
		stdprinter.Debug(
			"FileStore.CommitTransaction",
			"will index add",
			z.Path,
		)

		u.Index.Add(z)
	}

	for _, z := range u.ZettelsForActions(TransactionActionModified) {
		stdprinter.Debug(
			"FileStore.CommitTransaction",
			"will index modify",
			z.Path,
		)

		u.Index.Update(z)
	}

	for _, z := range u.ZettelsForActions(TransactionActionDeleted) {
		stdprinter.Debug(
			"FileStore.CommitTransaction",
			"will index delete",
			z.Path,
		)

		u.Index.Delete(z)
	}

	stdprinter.Debug("did commit transaction")

	err = u.CacheIndex()

	if err != nil {
		err = xerrors.Errorf("failed to cache index: %w", err)
		return
	}

	return
}
