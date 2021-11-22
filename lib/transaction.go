package lib

import (
	"sync"

	"github.com/friedenberg/z/lib/zettel"
	"github.com/friedenberg/z/util/stdprinter"
)

func MakeTransaction() (t *Transaction) {
	t = &Transaction{
		Locker:  &sync.Mutex{},
		actions: make([]TransactionEntry, 0),
	}

	return
}

//TODO-P1 handle cases where files or zettles are just opened but not edited
type Transaction struct {
	ShouldSkipCommit   bool
	ShouldCopyFiles    bool
	IsFinalTransaction bool
	sync.Locker
	actions []TransactionEntry
}

func (t *Transaction) Set(z *zettel.Zettel, action TransactionAction) {
	if z == nil {
		return
	}

	stdprinter.Debug(
		"Transaction.Set",
		"begin",
		z.Path,
		action,
	)

	t.Lock()
	defer t.Unlock()

	switch action {
	case TransactionActionNone:
		//TODO-P4
		panic("TODO")

	default:
		ne := TransactionEntry{
			Zettel:            z,
			TransactionAction: action,
		}

		t.actions = append(t.actions, ne)
	}

	stdprinter.Debug(
		"Transaction.Set",
		"end",
		z.Path,
		action,
	)
}

func (t Transaction) Len() int {
	t.Lock()
	defer t.Unlock()

	return len(t.actions)
}

func (t Transaction) ZettelsForActions(action TransactionAction) (zs ZettelSlice) {
	t.Lock()
	defer t.Unlock()

	zs = make([]*zettel.Zettel, 0, len(t.actions))

	for _, ze := range t.actions {
		stdprinter.Debug("ZettelsForActions")
		if ze.TransactionAction == action {
			zs = append(zs, ze.Zettel)
		}
	}

	stdprinter.Debugf("len: %d", len(zs))

	return
}

func (t Transaction) Paths() (f []string) {
	f = make([]string, 0, len(t.actions))

	for _, ze := range t.actions {
		f = append(f, ze.Path)
	}

	return
}
