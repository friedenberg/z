package lib

import (
	"sync"

	"github.com/friedenberg/z/lib/zettel"
)

func MakeTransaction() (t Transaction) {
	t = Transaction{
		Locker:  &sync.Mutex{},
		actions: make(map[zettel.Id]TransactionEntry),
	}

	return
}

//TODO-P1 handle cases where files or zettles are just opened but not edited
type Transaction struct {
	ShouldSkipCommit bool
	ShouldCopyFiles  bool
	sync.Locker
	actions map[zettel.Id]TransactionEntry
}

func (t Transaction) Set(z *Zettel, action TransactionAction) {
	if z == nil {
		return
	}

	t.Lock()
	defer t.Unlock()

	switch action {
	case TransactionActionNone:
		delete(t.actions, zettel.Id(z.Id))

	default:
		t.actions[zettel.Id(z.Id)] = TransactionEntry{
			Zettel:            z,
			TransactionAction: action,
		}
	}
}

func (t Transaction) Len() int {
	t.Lock()
	defer t.Unlock()

	return len(t.actions)
}

func (t Transaction) ZettelsForActions(action TransactionAction) (zs ZettelSlice) {
	t.Lock()
	defer t.Unlock()

	zs = make([]*Zettel, 0, len(t.actions))

	for _, ze := range t.actions {
		if ze.TransactionAction == action {
			zs = append(zs, ze.Zettel)
		}
	}

	return
}

func (t Transaction) Paths() (f []string) {
	f = make([]string, 0, len(t.actions))

	for _, ze := range t.actions {
		f = append(f, ze.Path)
	}

	return
}
