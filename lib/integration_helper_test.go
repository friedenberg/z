package lib

import (
	"io/ioutil"
	"reflect"
	"strings"
	"sync"
	"testing"

	"github.com/friedenberg/z/lib/zettel"
)

func printIndex(t *testing.T, u *Umwelt) {
	s, err := ioutil.ReadFile(u.GetIndexPath())

	if err != nil {
		t.Fatalf("failed to read index during debugging: %s", err)
	}

	t.Errorf("index:\n%s", s)
}

func assertNoTags(t *testing.T, u *Umwelt, z *zettel.Zettel, tag string) {
	t.Helper()

	i := u.Index
	id := zettel.Id(z.Id)

	assertIdsEmpty(t, z, i.Tags, tag)
	assertValuesEmpty(t, i.Tags, id, tag)
}

func assertTagInIndex(t *testing.T, u *Umwelt, z *zettel.Zettel, tag string) {
	t.Helper()
	i := u.Index
	ts := i.Tags
	id := zettel.Id(z.Id)

	if ss, ok := ts.GetValues(id); ok {
		if !ss.Contains(tag) {
			t.Errorf("missing expected tags in tag index: %s", tag)
		}
	} else {
		t.Errorf("missing expected zettel in tag index: %s", id)
	}

	if is, ok := ts.GetIds(tag); ok {
		if !is.Contains(id) {
			t.Errorf("missing expected zettel id in tag index id set: %s", id)
		}
	} else {
		t.Errorf("missing expected tag in tag index: %s", tag)
		printIndex(t, u)
	}
}

func assertZettelExistsInIndex(t *testing.T, u *Umwelt, z *zettel.Zettel) {
	i := u.Index
	zs := i.Zettels

	if _, ok := zs[zettel.Id(z.Id)]; !ok {
		t.Errorf("missing expected zettel in index: %s", zettel.Id(z.Id))
	}
}

func assertTransactionSuccessful(t *testing.T, u *Umwelt, count int) {
	t.Helper()

	transactionEntryCount := u.Transaction.Len()

	if transactionEntryCount != count {
		t.Fatalf("wrong number of added zettels: expected %d, got %d", count, transactionEntryCount)
	}

	err := u.Kasten.CommitTransaction(u)

	if err != nil {
		t.Fatalf("failed to commit transaction: %s", err)
	}

	err = u.CacheIndex()

	if err != nil {
		t.Fatalf("failed to write index: %s", err)
	}

	err = u.LoadIndexFromCache()

	if err != nil {
		t.Fatalf("failed to reload index: %s", err)
	}
}

func assertZettelMatchesFileSystem(t *testing.T, z *zettel.Zettel, expected string) {
	t.Helper()

	if z.Path == "" {
		t.Fatalf("zettel path is empty")
	}

	zettelFile := z.Path
	actual1, err := ioutil.ReadFile(zettelFile)

	if err != nil {
		t.Fatalf("failed to read written zettel (%s): %s", zettelFile, err)
	}

	actual := string(actual1)

	if actual != expected {
		t.Errorf("\nexpected:\n%#v\nactual:\n%#v", expected, actual)
	}
}

func makeZettel(t *testing.T, u *Umwelt, contents string) (z *zettel.Zettel) {
	t.Helper()

	z = &zettel.Zettel{}

	sr := strings.NewReader(contents)
	err := z.ReadFrom(sr, true)

	if err != nil {
		t.Fatalf("failed to hydrate zettel: %s", err)
	}

	u.Set(z, TransactionActionAdded)

	return
}

type TestFileStore struct {
	testing.T
	FileStore
	lastId int
	sync.Locker
}

func (k *TestFileStore) NewId() (i zettel.Id, err error) {
	k.Lock()
	defer k.Unlock()

	k.lastId += 1
	i = zettel.Id(k.lastId)

	return
}

// func (k *TestFileStore) CommitTransaction(u *Umwelt) (err error) {
// 	for _, z := range u.Transaction.actions {
// 		z.Path = zettel.Id(z.Id).String()
// 	}

// 	err = k.FileStore.CommitTransaction(u)

// 	return
// }

func makeUmwelt(t *testing.T) (u *Umwelt) {
	t.Helper()
	var err error

	u, err = MakeUmwelt(Config{}, t.TempDir())

	if err != nil {
		t.Fatalf("failed to make umwelt: %s", err)
	}

	u.Kasten = &TestFileStore{
		T:      *t,
		Locker: &sync.Mutex{},
		FileStore: FileStore{
			basePath: t.TempDir(),
		},
	}

	err = u.Kasten.Init(u, nil)

	if err != nil {
		t.Fatalf("failed to init kasten: %s", err)
	}

	return
}

func assertIdsEmpty(t *testing.T, z *zettel.Zettel, m zettel.MultiMap, k string) {
	t.Helper()
	actual1, ok := m.GetIds(k)

	if ok && actual1.Contains(zettel.Id(z.Id)) {
		t.Errorf("expected not to find %v -> %v", k, z.Id)
	}
}

func assertValuesEmpty(t *testing.T, m zettel.MultiMap, id zettel.Id, k string) {
	t.Helper()
	actual1, ok := m.GetValues(id)

	if ok && actual1.Contains(k) {
		t.Errorf("expected not to find %v -> %v", k, id)
	}
}

func assertValuesEqual(t *testing.T, m zettel.MultiMap, id zettel.Id, expected []string) {
	t.Helper()
	actual1, ok := m.GetValues(id)

	if !ok {
		t.Errorf("expected some values, but got none at all")
		return
	}

	actual := actual1.Slice()

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("expected %s but got %s", expected, actual)
	}
}

func assertIdsEqual(t *testing.T, m zettel.MultiMap, k string, expected []zettel.Id) {
	t.Helper()
	actual1, ok := m.GetIds(k)

	if !ok {
		t.Errorf("expected some values, but got none at all")
		return
	}

	actual := actual1.Slice()

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("expected %s but got %s", expected, actual)
	}
}
